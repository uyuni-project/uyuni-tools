// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

func installForKubernetes(
	globalFlags *types.GlobalFlags,
	flags *kubernetes.KubernetesServerFlags,
	cmd *cobra.Command,
	args []string,
) error {
	if _, err := exec.LookPath("kubectl"); err != nil {
		return fmt.Errorf(L("install kubectl before running this command"))
	}

	flags.Installation.CheckParameters(cmd, "kubectl")

	fqdn := args[0]
	if err := shared_utils.IsValidFQDN(fqdn); err != nil {
		return err
	}

	namespace := flags.Helm.Uyuni.Namespace
	// Create the namespace if not present
	if err := kubernetes.CreateNamespace(namespace); err != nil {
		return err
	}

	// TODO Is there a Hub API deployment?

	// TODO If there is already a server deployment, inspect it

	// TODO Run sanity checks for upgrade

	// TODO Get the fqdn from the inspected data if possible. Ignore difference witth input value for now.

	// TODO Scale down server deployment if present to upgrade the DB

	serverImage, err := shared_utils.ComputeImage(flags.Image.Registry, shared_utils.DefaultTag, flags.Image)
	if err != nil {
		return shared_utils.Errorf(err, L("failed to compute image URL"))
	}

	mounts := kubernetes.GetServerMounts()
	mounts = kubernetes.TuneMounts(mounts, &flags.Volumes)

	// TODO Only create PVCs if needed
	if err := shared_kubernetes.CreatePersistentVolumeClaims(namespace, mounts); err != nil {
		return err
	}

	// TODO Run the DB Upgrade job if needed

	// TODO Run DB finalization job

	// TODO Run Post upgrade job

	// Extract some data from the cluster to guess how to configure Uyuni.
	clusterInfos, err := shared_kubernetes.CheckCluster()
	if err != nil {
		return err
	}

	// Install the traefik / nginx config on the node
	// This will never be done in an operator.
	needsHub := flags.HubXmlrpc.Replicas > 0
	if err := kubernetes.DeployNodeConfig(
		namespace, clusterInfos, needsHub, flags.Installation.Debug.Java,
	); err != nil {
		return err
	}

	// Deploy the SSL CA and server certificates
	var caIssuer string
	sslFlags := flags.Installation.Ssl
	if sslFlags.UseExisting() {
		if err := kubernetes.DeployExistingCertificate(flags.Helm.Uyuni.Namespace, &sslFlags); err != nil {
			return err
		}
	} else {
		issuer, err := kubernetes.DeployGeneratedCa(
			&flags.Helm, &sslFlags, clusterInfos.GetKubeconfig(), fqdn, flags.Image.PullPolicy,
		)

		if err != nil {
			return shared_utils.Errorf(err, L("cannot deploy certificate"))
		}
		caIssuer = issuer
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

	// Start the server
	if err := kubernetes.CreateServerDeployment(
		namespace, serverImage, flags.Image.PullPolicy, flags.Installation.TZ, flags.Installation.Debug.Java,
		flags.Volumes.Mirror, pullSecret,
	); err != nil {
		return err
	}

	// Create the services
	if err := kubernetes.CreateServices(namespace, flags.Installation.Debug.Java); err != nil {
		return err
	}

	if clusterInfos.Ingress == "traefik" {
		// Create the Traefik routes
		if err := kubernetes.CreateTraefikRoutes(namespace, needsHub, flags.Installation.Debug.Java); err != nil {
			return err
		}
	}

	// Wait for the server deployment to have a running pod before trying to set it up.
	if err := shared_kubernetes.WaitForRunningDeployment(namespace, kubernetes.ServerDeployName); err != nil {
		return err
	}

	// TODO Run the setup only if it hasn't be done before: this is a one-off task.
	// Run the setup. This runs an exec into the running deployment.
	// TODO Ideally we would need a job running at an earlier stage to persist the logs in a kubernetes-friendly way.
	cnx := shared.NewConnection("kubectl", "", shared_kubernetes.ServerFilter)

	if err := adm_utils.RunSetup(
		cnx, &flags.ServerFlags, fqdn, map[string]string{"NO_SSL": "Y"},
	); err != nil {
		if stopErr := shared_kubernetes.Stop(namespace, shared_kubernetes.ServerApp); stopErr != nil {
			log.Error().Msgf(L("Failed to stop service: %v"), stopErr)
		}
		return err
	}

	// Store the DB credentials in a secret.
	dbFlags := flags.Installation.Db
	if err := kubernetes.CreateDbSecret(namespace, kubernetes.DbSecret, dbFlags.User, dbFlags.Password); err != nil {
		return err
	}

	deploymentsStarting := []string{}

	// Start the Coco Deployments if requested.
	if flags.Coco.Replicas > 0 {
		cocoImage, err := shared_utils.ComputeImage(flags.Image.Registry, flags.Image.Tag, flags.Coco.Image)
		if err != nil {
			return err
		}
		if err := kubernetes.StartCocoDeployment(
			namespace, cocoImage, flags.Image.PullPolicy, flags.Coco.Replicas,
			dbFlags.Port, dbFlags.Name,
		); err != nil {
			return err
		}
		deploymentsStarting = append(deploymentsStarting, kubernetes.CocoDeployName)
	}

	// In an operator mind, the user would just change the custom resource to enable the feature.
	if needsHub {
		// Install Hub API deployment, service
		hubApiImage, err := shared_utils.ComputeImage(flags.Image.Registry, flags.Image.Tag, flags.HubXmlrpc.Image)
		if err != nil {
			return err
		}
		if err := kubernetes.InstallHubApi(namespace, hubApiImage, flags.Image.PullPolicy); err != nil {
			return err
		}
		deploymentsStarting = append(deploymentsStarting, kubernetes.HubApiDeployName)
	}

	// Wait for all the other deployments to be ready
	if err := shared_kubernetes.WaitForDeployments(namespace, deploymentsStarting...); err != nil {
		return err
	}

	return nil
}
