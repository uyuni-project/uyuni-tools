// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"encoding/base64"
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Reconcile upgrades, migrate or install the server.
func Reconcile(flags *KubernetesServerFlags, fqdn string) error {
	if _, err := exec.LookPath("kubectl"); err != nil {
		return fmt.Errorf(L("install kubectl before running this command"))
	}

	if err := utils.IsValidFQDN(fqdn); err != nil {
		return err
	}

	namespace := flags.Helm.Uyuni.Namespace
	// Create the namespace if not present
	if err := CreateNamespace(namespace); err != nil {
		return err
	}

	serverImage, err := utils.ComputeImage(flags.Image.Registry, utils.DefaultTag, flags.Image)
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	// TODO Is there a Hub API deployment?

	cnx := shared.NewConnection("kubectl", "", kubernetes.ServerFilter)

	// Do we have an existing deployment to upgrade?
	// This can be freshly synchronized data from a migration or a running instance to upgrade.
	hasDatabase := kubernetes.HasVolume(namespace, "var-pgsql")
	var inspectedData utils.ServerInspectData
	if hasDatabase {
		// Inspect the image and the existing volumes
		data, err := inspectServer(namespace, serverImage, flags.Image.PullPolicy)
		if err != nil {
			return err
		}
		inspectedData = *data

		// Do we have a running server deploy? which version is it?
		// If there is no deployment / image, don't check the uyuni / SUMA upgrades
		var runningData *utils.ServerInspectData
		if runningImage := getRunningServerImage(namespace); runningImage != "" {
			runningData, err = inspectServer(namespace, runningImage, "Never")
			if err != nil {
				return err
			}
		}

		// Run sanity checks for upgrade
		if err := adm_utils.SanityCheck(runningData, &inspectedData, serverImage); err != nil {
			return err
		}

		// Get the fqdn from the inspected data if possible. Ignore difference with input value for now.
		fqdn = inspectedData.Fqdn

		if inspectedData.CurrentPgVersion != inspectedData.ImagePgVersion {
			// Scale down server deployment if present to upgrade the DB
			if err := kubernetes.ReplicasTo(namespace, kubernetes.ServerApp, 0); err != nil {
				return utils.Errorf(err, L("cannot set replica to 0"))
			}
		}
	}

	mounts := GetServerMounts()
	mounts = TuneMounts(mounts, &flags.Volumes)

	if err := kubernetes.CreatePersistentVolumeClaims(namespace, mounts); err != nil {
		return err
	}

	oldPgVersion := inspectedData.CurrentPgVersion
	newPgVersion := inspectedData.ImagePgVersion

	if hasDatabase {
		// Run the DB Upgrade job if needed
		if oldPgVersion < newPgVersion {
			if err := StartDbUpgradeJob(
				namespace, flags.Image.Registry, flags.Image, flags.DbUpgradeImage,
				oldPgVersion, newPgVersion,
			); err != nil {
				return err
			}

			// Wait for ever for the job to finish: the duration of this job depends on the amount of data to upgrade
			if err := kubernetes.WaitForJob(namespace, DbUpgradeJobName, -1); err != nil {
				return err
			}
		} else if oldPgVersion > newPgVersion {
			return fmt.Errorf(
				L("downgrading database from PostgreSQL %[1]d to %[2]d is not supported"), oldPgVersion, newPgVersion)
		}

		// Run DB finalization job
		schemaUpdateRequired := oldPgVersion != newPgVersion
		if err := StartDbFinalizeJob(
			namespace, serverImage, flags.Image.PullPolicy, schemaUpdateRequired, true,
		); err != nil {
			return err
		}

		// Wait for ever for the job to finish: the duration of this job depends on the amount of data to reindex
		if err := kubernetes.WaitForJob(namespace, DbFinalizeJobName, -1); err != nil {
			return err
		}

		// Run the Post Upgrade job
		if err := StartPostUpgradeJob(namespace, serverImage, flags.Image.PullPolicy); err != nil {
			return err
		}

		if err := kubernetes.WaitForJob(namespace, PostUpgradeJobName, 60); err != nil {
			return err
		}
	}

	// Extract some data from the cluster to guess how to configure Uyuni.
	clusterInfos, err := kubernetes.CheckCluster()
	if err != nil {
		return err
	}

	// Install the traefik / nginx config on the node
	// This will never be done in an operator.
	needsHub := flags.HubXmlrpc.Replicas > 0
	if err := DeployNodeConfig(namespace, clusterInfos, needsHub, flags.Installation.Debug.Java); err != nil {
		return err
	}

	// Deploy the SSL CA and server certificates
	var caIssuer string
	if flags.Installation.Ssl.UseExisting() {
		if err := DeployExistingCertificate(flags.Helm.Uyuni.Namespace, &flags.Installation.Ssl); err != nil {
			return err
		}
	} else {
		// cert-manager is not required for 3rd party certificates, only if we have the CA key.
		// Note that in an operator we won't be able to install cert-manager and just wait for it to be installed.
		kubeconfig := clusterInfos.GetKubeconfig()

		if err := InstallCertManager(&flags.Helm, kubeconfig, flags.Image.PullPolicy); err != nil {
			return utils.Errorf(err, L("cannot install cert manager"))
		}

		if flags.Installation.Ssl.UseMigratedCa() {
			// Convert CA to RSA to use in a Kubernetes TLS secret.
			// In an operator we would have to fail now if there is no SSL password as we cannot prompt it.
			ca := ssl.SslPair{
				Key: base64.StdEncoding.EncodeToString(
					ssl.GetRsaKey(flags.Installation.Ssl.Ca.Key, flags.Installation.Ssl.Password),
				),
				Cert: base64.StdEncoding.EncodeToString(ssl.StripTextFromCertificate(flags.Installation.Ssl.Ca.Root)),
			}

			// Install the cert-manager issuers
			if _, err := DeployReusedCa(namespace, &ca); err != nil {
				return err
			}
		} else {
			if err := DeployGeneratedCa(flags.Helm.Uyuni.Namespace, &flags.Installation.Ssl, fqdn); err != nil {
				return err
			}
		}

		// Wait for issuer to be ready
		if err := waitForIssuer(flags.Helm.Uyuni.Namespace, CaIssuerName); err != nil {
			return err
		} else {
		}

		// Extract the CA cert into uyuni-ca config map as the container shouldn't have the CA secret
		if err := extractCaCertToConfig(flags.Helm.Uyuni.Namespace); err != nil {
			return err
		}
		caIssuer = CaIssuerName
	}

	// Create the Ingress routes before the deployments as those are triggering
	// the creation of the uyuni-cert secret from cert-manager.
	if err := CreateIngress(namespace, fqdn, caIssuer, clusterInfos.Ingress); err != nil {
		return err
	}

	// Wait for uyuni-cert secret to be ready
	kubernetes.WaitForSecret(namespace, CertSecretName)

	// Create a secret using SCC credentials if any are provided
	pullSecret, err := kubernetes.GetSccSecret(flags.Helm.Uyuni.Namespace, &flags.Installation.Scc)
	if err != nil {
		return err
	}

	// Start the server
	if err := CreateServerDeployment(
		namespace, serverImage, flags.Image.PullPolicy, flags.Installation.TZ, flags.Installation.Debug.Java,
		flags.Volumes.Mirror, pullSecret,
	); err != nil {
		return err
	}

	// Create the services
	if err := CreateServices(namespace, flags.Installation.Debug.Java); err != nil {
		return err
	}

	if clusterInfos.Ingress == "traefik" {
		// Create the Traefik routes
		if err := CreateTraefikRoutes(namespace, needsHub, flags.Installation.Debug.Java); err != nil {
			return err
		}
	}

	// Wait for the server deployment to have a running pod before trying to set it up.
	if err := kubernetes.WaitForRunningDeployment(namespace, ServerDeployName); err != nil {
		return err
	}

	// TODO Run the setup only if it hasn't be done before: this is a one-off task.
	// Run the setup. This runs an exec into the running deployment.
	// TODO Ideally we would need a job running at an earlier stage to persist the logs in a kubernetes-friendly way.

	if err := adm_utils.RunSetup(
		cnx, &flags.ServerFlags, fqdn, map[string]string{"NO_SSL": "Y"},
	); err != nil {
		if stopErr := kubernetes.Stop(namespace, kubernetes.ServerApp); stopErr != nil {
			log.Error().Msgf(L("Failed to stop service: %v"), stopErr)
		}
		return err
	}

	// Store the DB credentials in a secret.
	if flags.Installation.Db.User != "" && flags.Installation.Db.Password != "" {
		if err := CreateDbSecret(
			namespace, DbSecret, flags.Installation.Db.User, flags.Installation.Db.Password,
		); err != nil {
			return err
		}
	}

	deploymentsStarting := []string{}

	// Start the Coco Deployments if requested.
	if flags.Coco.Replicas > 0 {
		cocoImage, err := utils.ComputeImage(flags.Image.Registry, flags.Image.Tag, flags.Coco.Image)
		if err != nil {
			return err
		}
		if err := StartCocoDeployment(
			namespace, cocoImage, flags.Image.PullPolicy, flags.Coco.Replicas,
			flags.Installation.Db.Port, flags.Installation.Db.Name,
		); err != nil {
			return err
		}
		deploymentsStarting = append(deploymentsStarting, CocoDeployName)
	}

	// In an operator mind, the user would just change the custom resource to enable the feature.
	if needsHub {
		// Install Hub API deployment, service
		hubApiImage, err := utils.ComputeImage(flags.Image.Registry, flags.Image.Tag, flags.HubXmlrpc.Image)
		if err != nil {
			return err
		}
		if err := InstallHubApi(namespace, hubApiImage, flags.Image.PullPolicy); err != nil {
			return err
		}
		deploymentsStarting = append(deploymentsStarting, HubApiDeployName)
	}

	// Wait for all the other deployments to be ready
	if err := kubernetes.WaitForDeployments(namespace, deploymentsStarting...); err != nil {
		return err
	}

	return nil
}
