// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const migrationDataPvcName = "migration-data"

func migrateToKubernetes(
	globalFlags *types.GlobalFlags,
	flags *kubernetes.KubernetesServerFlags,
	cmd *cobra.Command,
	args []string,
) error {
	namespace := flags.Helm.Uyuni.Namespace
	// Create the namespace if not present
	if err := kubernetes.CreateNamespace(namespace); err != nil {
		return err
	}

	// Check the for the required SSH key and configuration
	if err := checkSsh(namespace, &flags.Ssh); err != nil {
		return err
	}

	serverImage, err := utils.ComputeImage(flags.Image.Registry, utils.DefaultTag, flags.Image)
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	fqdn := args[0]
	if err := utils.IsValidFQDN(fqdn); err != nil {
		return err
	}

	mounts := kubernetes.GetServerMounts()
	mounts = kubernetes.TuneMounts(mounts, &flags.Volumes)

	// Add a mount and volume for the extracted data
	migrationDataVolume := types.VolumeMount{Name: migrationDataPvcName, MountPath: "/var/lib/uyuni-tools"}
	migrationMounts := append(mounts, migrationDataVolume)

	if err := shared_kubernetes.CreatePersistentVolumeClaims(namespace, migrationMounts); err != nil {
		return err
	}

	if err = startMigrationJob(
		namespace,
		serverImage,
		flags.Image.PullPolicy,
		fqdn,
		flags.Migration.User,
		flags.Migration.Prepare,
		migrationMounts,
	); err != nil {
		return err
	}

	// Wait for ever for the job to finish: the duration of this job depends on the amount of data to copy
	if err := shared_kubernetes.WaitForJob(namespace, migrationJobName, -1); err != nil {
		return err
	}

	// Read the extracted data from the migration volume
	extractedData, err := extractMigrationData(namespace, serverImage, migrationDataVolume)
	if err != nil {
		return err
	}

	oldPgVersion := extractedData.Data.CurrentPgVersion
	newPgVersion := extractedData.Data.ImagePgVersion

	// Run the DB Migration job if needed
	if oldPgVersion < newPgVersion {
		if err := kubernetes.StartDbUpgradeJob(
			namespace, flags.Image.Registry, flags.Image, flags.DbUpgradeImage,
			oldPgVersion, newPgVersion,
		); err != nil {
			return err
		}

		// Wait for ever for the job to finish: the duration of this job depends on the amount of data to upgrade
		if err := shared_kubernetes.WaitForJob(namespace, kubernetes.DbUpgradeJobName, -1); err != nil {
			return err
		}
	} else if oldPgVersion > newPgVersion {
		return fmt.Errorf(
			L("downgrading database from PostgreSQL %[1]d to %[2]d is not supported"), oldPgVersion, newPgVersion)
	}

	// Run the DB Finalization job
	schemaUpdateRequired := oldPgVersion != newPgVersion
	if err := kubernetes.StartDbFinalizeJob(
		namespace, serverImage, flags.Image.PullPolicy, schemaUpdateRequired, true,
	); err != nil {
		return err
	}

	// Wait for ever for the job to finish: the duration of this job depends on the amount of data to reindex
	if err := shared_kubernetes.WaitForJob(namespace, kubernetes.DbFinalizeJobName, -1); err != nil {
		return err
	}

	// Run the Post Upgrade job
	if err := kubernetes.StartPostUpgradeJob(namespace, serverImage, flags.Image.PullPolicy); err != nil {
		return err
	}

	if err := shared_kubernetes.WaitForJob(namespace, kubernetes.PostUpgradeJobName, 60); err != nil {
		return err
	}

	// Extract some data from the cluster to guess how to configure Uyuni.
	clusterInfos, err := shared_kubernetes.CheckCluster()
	if err != nil {
		return err
	}

	// Install the traefik / nginx config on the node
	// This will never be done in an operator.
	needsHub := flags.HubXmlrpc.Replicas > 0
	if err := kubernetes.DeployNodeConfig(namespace, clusterInfos, needsHub, extractedData.Data.Debug); err != nil {
		return err
	}

	// Deploy the SSL CA and server certificates
	var caIssuer string
	if extractedData.CaKey != "" {
		// cert-manager is not required for 3rd party certificates, only if we have the CA key.
		// Note that in an operator we won't be able to install cert-manager and just wait for it to be installed.
		kubeconfig := clusterInfos.GetKubeconfig()

		if err := kubernetes.InstallCertManager(&flags.Helm, kubeconfig, flags.Image.PullPolicy); err != nil {
			return utils.Errorf(err, L("cannot install cert manager"))
		}

		// Convert CA to RSA to use in a Kubernetes TLS secret.
		// In an operator we would have to fail now if there is no SSL password as we cannot prompt it.
		ca := ssl.SslPair{
			Key: base64.StdEncoding.EncodeToString(
				ssl.GetRsaKey(extractedData.CaKey, flags.Installation.Ssl.Password),
			),
			Cert: base64.StdEncoding.EncodeToString(ssl.StripTextFromCertificate(extractedData.CaCert)),
		}

		// Install the cert-manager issuers
		if _, err := kubernetes.DeployReusedCa(namespace, &ca); err != nil {
			return err
		}
		caIssuer = kubernetes.CaIssuerName
	} else {
		// Most likely a 3rd party certificate: cert-manager is not needed in this case
		if err := installExistingCertificate(namespace, extractedData); err != nil {
			return err
		}
	}

	// Create the Ingress routes before the deployments as those are triggering
	// the creation of the uyuni-cert secret from cert-manager.
	if err := kubernetes.CreateIngress(namespace, fqdn, caIssuer, clusterInfos.Ingress); err != nil {
		return err
	}

	// Wait for uyuni-cert secret to be ready
	shared_kubernetes.WaitForSecret(namespace, kubernetes.CertSecretName)

	// Create a secret using SCC credentials if any are provided
	pullSecret, err := shared_kubernetes.GetSccSecret(flags.Helm.Uyuni.Namespace, &flags.Installation.Scc)
	if err != nil {
		return err
	}

	deploymentsStarting := []string{kubernetes.ServerDeployName}
	// Start the server
	if err := kubernetes.CreateServerDeployment(
		namespace, serverImage, flags.Image.PullPolicy, extractedData.Data.Timezone, extractedData.Data.Debug,
		flags.Volumes.Mirror, pullSecret,
	); err != nil {
		return err
	}

	// Create the services
	if err := kubernetes.CreateServices(namespace, extractedData.Data.Debug); err != nil {
		return err
	}

	if clusterInfos.Ingress == "traefik" {
		// Create the Traefik routes
		if err := kubernetes.CreateTraefikRoutes(namespace, needsHub, extractedData.Data.Debug); err != nil {
			return err
		}
	}

	// Store the extracted DB credentials in a secret.
	if err := kubernetes.CreateDbSecret(
		namespace, kubernetes.DbSecret, extractedData.Data.DbUser, extractedData.Data.DbPassword,
	); err != nil {
		return err
	}

	// Start the Coco Deployments if requested.
	if flags.Coco.Replicas > 0 {
		cocoImage, err := utils.ComputeImage(flags.Image.Registry, flags.Image.Tag, flags.Coco.Image)
		if err != nil {
			return err
		}
		if err := kubernetes.StartCocoDeployment(
			namespace, cocoImage, flags.Image.PullPolicy, flags.Coco.Replicas,
			extractedData.Data.DbPort, extractedData.Data.DbName,
		); err != nil {
			return err
		}
		deploymentsStarting = append(deploymentsStarting, kubernetes.CocoDeployName)
	}

	// In an operator mind, the user would just change the custom resource to enable the feature.
	if extractedData.Data.HasHubXmlrpcApi {
		// Install Hub API deployment, service
		hubApiImage, err := utils.ComputeImage(flags.Image.Registry, flags.Image.Tag, flags.HubXmlrpc.Image)
		if err != nil {
			return err
		}
		if err := kubernetes.InstallHubApi(namespace, hubApiImage, flags.Image.PullPolicy); err != nil {
			return err
		}
		deploymentsStarting = append(deploymentsStarting, kubernetes.HubApiDeployName)
	}

	// Wait for all the deployments to be ready
	if err := shared_kubernetes.WaitForDeployments(namespace, deploymentsStarting...); err != nil {
		return err
	}

	return nil
}
