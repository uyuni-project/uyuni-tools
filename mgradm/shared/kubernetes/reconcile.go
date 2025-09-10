// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"

	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Reconcile upgrades, migrate or install the server.
func Reconcile(flags *KubernetesServerFlags, fqdn string) error {
	if _, err := exec.LookPath("kubectl"); err != nil {
		return errors.New(L("install kubectl before running this command"))
	}

	namespace := flags.Kubernetes.Uyuni.Namespace
	// Create the namespace if not present
	if err := CreateNamespace(namespace); err != nil {
		return err
	}

	serverImage, err := utils.ComputeImage(flags.Image.Registry.Host, utils.DefaultTag, flags.Image)
	if err != nil {
		return utils.Error(err, L("failed to compute image URL"))
	}

	// Create a secret using SCC credentials if any are provided
	pullSecret, err := kubernetes.GetRegistrySecret(
		flags.Kubernetes.Uyuni.Namespace, &flags.Image.Registry, kubernetes.ServerApp,
	)
	if err != nil {
		return err
	}

	// Do we have an existing deployment to upgrade?
	// This can be freshly synchronized data from a migration or a running instance to upgrade.
	hasDeployment := kubernetes.HasDeployment(namespace, kubernetes.ServerFilter)

	// Check that the postgresql PVC is bound to a Volume.
	hasDatabase := kubernetes.HasVolume(namespace, "var-pgsql")
	isMigration := hasDatabase && !hasDeployment

	cocoReplicas := kubernetes.GetReplicas(namespace, CocoDeployName)
	if cocoReplicas != 0 && !flags.Coco.IsChanged {
		// Upgrade: detect the number of running coco replicas
		flags.Coco.Replicas = cocoReplicas
	}

	var inspectedData utils.ServerInspectData
	if hasDatabase {
		// Inspect the image and the existing volumes
		data, err := kubernetes.InspectServer(namespace, serverImage, flags.Image.PullPolicy, pullSecret)
		if err != nil {
			return err
		}
		inspectedData = *data

		// Use the inspected DB port and name if not defined in the flags
		if flags.Installation.DB.Port == 0 && data.DBPort != 0 {
			flags.Installation.DB.Port = data.DBPort
		}

		if flags.Installation.DB.Name == "" && data.DBName != "" {
			flags.Installation.DB.Name = data.DBName
		}

		// TODO Do we have a running database deployment?

		// Do we have a running server deploy? which version is it?
		// If there is no deployment / image, don't check the uyuni / SUMA upgrades
		// TODO If the DB is already in a separate deployment, there surely is no need to run an inspect on the server pod.
		var runningData *utils.ServerInspectData
		if runningImage := getRunningServerImage(namespace); runningImage != "" {
			runningData, err = kubernetes.InspectServer(namespace, runningImage, "Never", pullSecret)
			if err != nil {
				return err
			}
		}

		// Run sanity checks for upgrade
		if err := adm_utils.SanityCheck(runningData, &inspectedData); err != nil {
			return err
		}

		// Get the fqdn from the inspected data if possible. Ignore difference with input value for now.
		fqdn = inspectedData.Fqdn

		if hasDeployment {
			// Scale down all deployments relying on the DB since it will be brought down during upgrade.
			if cocoReplicas > 0 {
				if err := kubernetes.ReplicasTo(namespace, CocoDeployName, 0); err != nil {
					return utils.Error(err, L("cannot set confidential computing containers replicas to 0"))
				}
			}

			// Scale down server deployment if present to upgrade the DB
			if err := kubernetes.ReplicasTo(namespace, ServerDeployName, 0); err != nil {
				return utils.Error(err, L("cannot set server replicas to 0"))
			}

			// TODO Scale down the DB container?
		}
	}

	// Don't check the FQDN too early or we may not have it in case of upgrade.
	if err := utils.IsValidFQDN(fqdn); err != nil {
		return err
	}

	// TODO IsLocal() is not enough for Kubernetes as users can define their own db / reportdb service pointing
	// to whatever they want
	localDB := flags.Installation.DB.IsLocal()

	mounts := GetServerMounts()
	if localDB {
		mounts = append(mounts, utils.VarPgsqlDataVolumeMount)
	}
	mounts = TuneMounts(mounts, &flags.Volumes)

	if err := kubernetes.CreatePersistentVolumeClaims(namespace, mounts); err != nil {
		return err
	}

	if hasDatabase {
		oldPgVersion := inspectedData.CommonInspectData.CurrentPgVersion
		newPgVersion := inspectedData.DBInspectData.ImagePgVersion

		// TODO Split DB upgrade if needed or merge it in another job (which?)

		// Run the DB Upgrade job if needed
		if oldPgVersion < newPgVersion {
			jobName, err := StartDBUpgradeJob(
				namespace, flags.Image, flags.DBUpgradeImage, pullSecret,
				oldPgVersion, newPgVersion,
			)
			if err != nil {
				return err
			}

			// Wait for ever for the job to finish: the duration of this job depends on the amount of data to upgrade
			if err := kubernetes.WaitForJob(namespace, jobName, -1); err != nil {
				return err
			}
		} else if oldPgVersion > newPgVersion {
			return fmt.Errorf(
				L("downgrading database from PostgreSQL %[1]d to %[2]d is not supported"), oldPgVersion, newPgVersion)
		}

		// Run DB finalization job
		schemaUpdateRequired := oldPgVersion != newPgVersion
		jobName, err := StartDBFinalizeJob(
			namespace, serverImage, flags.Image.PullPolicy, pullSecret, schemaUpdateRequired, isMigration,
		)
		if err != nil {
			return err
		}

		// Wait for ever for the job to finish: the duration of this job depends on the amount of data to reindex
		if err := kubernetes.WaitForJob(namespace, jobName, -1); err != nil {
			return err
		}

		// Run the Post Upgrade job
		jobName, err = StartPostUpgradeJob(namespace, serverImage, flags.Image.PullPolicy, pullSecret)
		if err != nil {
			return err
		}

		if err := kubernetes.WaitForJob(namespace, jobName, 60); err != nil {
			return err
		}
	}

	// Extract some data from the cluster to guess how to configure Uyuni.
	clusterInfos, err := kubernetes.CheckCluster()
	if err != nil {
		return err
	}

	if replicas := kubernetes.GetReplicas(namespace, ServerDeployName); replicas > 0 && !flags.HubXmlrpc.IsChanged {
		// Upgrade: detect the number of existing hub xmlrpc replicas
		flags.HubXmlrpc.Replicas = replicas
	}
	needsHub := flags.HubXmlrpc.Replicas > 0

	// Install the traefik / nginx config on the node
	// This will never be done in an operator.
	if err := deployNodeConfig(namespace, clusterInfos, needsHub, flags.Installation.Debug.Java); err != nil {
		return err
	}

	// Deploy the SSL CA and server certificates
	var caIssuer string
	if flags.Installation.SSL.UseProvided() {
		if err := DeployExistingCertificate(flags.Kubernetes.Uyuni.Namespace, &flags.Installation.SSL); err != nil {
			return err
		}
	} else if !HasIssuer(namespace, kubernetes.CAIssuerName) {
		// cert-manager is not required for 3rd party certificates, only if we have the CA key.
		// Note that in an operator we won't be able to install cert-manager and just wait for it to be installed.
		kubeconfig := clusterInfos.GetKubeconfig()

		if err := InstallCertManager(&flags.Kubernetes, kubeconfig, flags.Image.PullPolicy); err != nil {
			return utils.Error(err, L("cannot install cert manager"))
		}

		if flags.Installation.SSL.UseMigratedCa() {
			// Convert CA to RSA to use in a Kubernetes TLS secret.
			// In an operator we would have to fail now if there is no SSL password as we cannot prompt it.
			rootCA, err := os.ReadFile(flags.Installation.SSL.Ca.Root)
			if err != nil {
				return utils.Error(err, L("failed to read Root CA file"))
			}
			ca := types.SSLPair{
				Key: base64.StdEncoding.EncodeToString(
					ssl.GetRsaKey(flags.Installation.SSL.Ca.Key, flags.Installation.SSL.Password),
				),
				Cert: base64.StdEncoding.EncodeToString(ssl.StripTextFromCertificate(string(rootCA))),
			}

			// Install the cert-manager issuers
			if err := DeployReusedCA(namespace, &ca, fqdn); err != nil {
				return err
			}
		} else {
			if err := DeployGeneratedCA(flags.Kubernetes.Uyuni.Namespace, &flags.Installation.SSL, fqdn); err != nil {
				return err
			}
		}

		// Wait for issuer to be ready
		if err := waitForIssuer(flags.Kubernetes.Uyuni.Namespace, kubernetes.CAIssuerName); err != nil {
			return err
		}

		// Extract the CA cert into uyuni-ca config map as the container shouldn't have the CA secret
		if err := extractCACertToConfig(flags.Kubernetes.Uyuni.Namespace); err != nil {
			return err
		}
		caIssuer = kubernetes.CAIssuerName
	}

	// Create the Ingress routes before the deployments as those are triggering
	// the creation of the uyuni-cert secret from cert-manager.
	if err := CreateIngress(namespace, fqdn, caIssuer, clusterInfos.Ingress); err != nil {
		return err
	}

	// Wait for uyuni-cert secret to be ready
	kubernetes.WaitForSecret(namespace, kubernetes.CertSecretName)

	// Create the services
	if err := CreateServices(namespace, flags.Installation.Debug.Java); err != nil {
		return err
	}

	// Store the DB credentials in a secret.
	if flags.Installation.DB.User != "" && flags.Installation.DB.Password != "" {
		if err := CreateBasicAuthSecret(
			namespace, DBSecret, flags.Installation.DB.User, flags.Installation.DB.Password,
		); err != nil {
			return err
		}
	}

	if flags.Installation.ReportDB.User != "" && flags.Installation.ReportDB.Password != "" {
		if err := CreateBasicAuthSecret(
			namespace, ReportdbSecret, flags.Installation.ReportDB.User, flags.Installation.ReportDB.Password,
		); err != nil {
			return err
		}
	}

	if !hasDatabase {
		// Wait for the DB secrets: TLS, ReportDB and DB credentials
		kubernetes.WaitForSecret(namespace, DBSecret)
		kubernetes.WaitForSecret(namespace, ReportdbSecret)
		kubernetes.WaitForSecret(namespace, kubernetes.DBCertSecretName)

		if localDB {
			// Create the secret for admin credentials
			if err := CreateBasicAuthSecret(
				namespace, DBAdminSecret, flags.Installation.DB.Admin.User, flags.Installation.DB.Admin.Password,
			); err != nil {
				return err
			}
			kubernetes.WaitForSecret(namespace, DBAdminSecret)

			dbImage, err := utils.ComputeImage(flags.Image.Registry.Host, utils.DefaultTag, flags.Pgsql.Image)
			if err != nil {
				return utils.Error(err, L("failed to compute image URL"))
			}

			// Create the split DB deployment
			if err := CreateDBDeployment(
				namespace, dbImage, flags.Image.PullPolicy, pullSecret, flags.Installation.TZ,
			); err != nil {
				return err
			}
		}
	}

	// This SCCSecret is used to mount the env variable in the setup job and is different from the
	// pullSecret as it is of a different type: basic-auth vs docker.
	if flags.Installation.SCC.User != "" && flags.Installation.SCC.Password != "" {
		if err := CreateBasicAuthSecret(
			namespace, SCCSecret, flags.Installation.SCC.User, flags.Installation.SCC.Password,
		); err != nil {
			return err
		}
	}

	adminSecret := "admin-credentials"
	if flags.Installation.Admin.Login != "" && flags.Installation.Admin.Password != "" {
		if err := CreateBasicAuthSecret(
			namespace, adminSecret, flags.Installation.Admin.Login, flags.Installation.Admin.Password,
		); err != nil {
			return err
		}
	}

	// TODO For a migration or an upgrade this needs to be skipped
	// Run the setup script.
	// The script will be skipped if the server has already been setup.
	jobName, err := StartSetupJob(
		namespace, serverImage, kubernetes.GetPullPolicy(flags.Image.PullPolicy), pullSecret,
		flags.Volumes.Mirror, &flags.Installation, fqdn, adminSecret, DBSecret, ReportdbSecret, SCCSecret,
	)
	if err != nil {
		return err
	}

	if err := kubernetes.WaitForJob(namespace, jobName, 120); err != nil {
		return err
	}

	if clusterInfos.Ingress == "traefik" {
		// Create the Traefik routes
		if err := CreateTraefikRoutes(namespace, needsHub, flags.Installation.Debug.Java); err != nil {
			return err
		}
	}

	// Start the server
	if err := CreateServerDeployment(
		namespace, serverImage, flags.Image.PullPolicy, flags.Installation.TZ, flags.Installation.Debug.Java,
		flags.Volumes.Mirror, pullSecret,
	); err != nil {
		return err
	}

	deploymentsStarting := []string{ServerDeployName}

	// Start the Coco Deployments if requested.
	if replicas := kubernetes.GetReplicas(namespace, CocoDeployName); replicas != 0 && !flags.Coco.IsChanged {
		// Upgrade: detect the number of running coco replicas
		flags.Coco.Replicas = replicas
	}
	if flags.Coco.Replicas > 0 {
		cocoImage, err := utils.ComputeImage(flags.Image.Registry.Host, flags.Image.Tag, flags.Coco.Image)
		if err != nil {
			return err
		}
		if err := StartCocoDeployment(
			namespace, cocoImage, flags.Image.PullPolicy, pullSecret, flags.Coco.Replicas,
			flags.Installation.DB.Port, flags.Installation.DB.Name,
		); err != nil {
			return err
		}
		deploymentsStarting = append(deploymentsStarting, CocoDeployName)
	}

	// In an operator mind, the user would just change the custom resource to enable the feature.
	if needsHub {
		// Install Hub API deployment, service
		hubAPIImage, err := utils.ComputeImage(flags.Image.Registry.Host, flags.Image.Tag, flags.HubXmlrpc.Image)
		if err != nil {
			return err
		}
		if err := InstallHubAPI(namespace, hubAPIImage, flags.Image.PullPolicy, pullSecret); err != nil {
			return err
		}
		deploymentsStarting = append(deploymentsStarting, HubAPIDeployName)
	}

	// Wait for all the other deployments to be ready
	if err := kubernetes.WaitForDeployments(namespace, deploymentsStarting...); err != nil {
		return err
	}

	return nil
}
