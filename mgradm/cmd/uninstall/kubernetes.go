// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package uninstall

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const kubernetesBuilt = true

func uninstallForKubernetes(dryRun bool) {
	clusterInfos := shared_kubernetes.CheckCluster()
	kubeconfig := clusterInfos.GetKubeconfig()

	// TODO Find all the PVs related to the server if we want to delete them

	// Uninstall uyuni
	namespace := helmUninstall(kubeconfig, "uyuni", "", dryRun)

	// Remove the remaining configmap and secrets
	if namespace != "" {

		_, err := utils.RunCmdOutput(zerolog.TraceLevel, "kubectl", "-n", namespace, "get", "secret", "uyuni-ca")
		caSecret := "uyuni-ca"
		if err != nil {
			caSecret = ""
		}

		if dryRun {
			log.Info().Msgf("Would run kubectl delete -n %s configmap uyuni-ca", namespace)
			log.Info().Msgf("Would run kubectl delete -n %s secret uyuni-cert %s", namespace, caSecret)
		} else {
			log.Info().Msgf("Running kubectl delete -n %s configmap uyuni-ca", namespace)
			if err := utils.RunCmd("kubectl", "delete", "-n", namespace, "configmap", "uyuni-ca"); err != nil {
				log.Info().Err(err).Msgf("Failed deleting config map")
			}

			log.Info().Msgf("Running kubectl delete -n %s secret uyuni-cert %s", namespace, caSecret)

			args := []string{"delete", "-n", namespace, "secret", "uyuni-cert"}
			if caSecret != "" {
				args = append(args, caSecret)
			}
			err := utils.RunCmd("kubectl", args...)
			if err != nil {
				log.Info().Err(err).Msgf("Failed deleting secret")
			}
		}
	}

	// TODO Remove the PVs or wait for their automatic removal if purge is requested
	// Also wait if the PVs are dynamic with Delete reclaim policy but the user didn't ask to purge them
	// Since some storage plugins don't handle Delete policy, we may need to check for error events to avoid infinite loop

	// Uninstall cert-manager if we installed it
	helmUninstall(kubeconfig, "cert-manager", "-linstalledby=mgradm", dryRun)

	// Remove the K3s Traefik config
	if clusterInfos.IsK3s() {
		shared_kubernetes.UninstallK3sTraefikConfig(dryRun)
	}

	// Remove the rke2 nginx config
	if clusterInfos.IsRke2() {
		shared_kubernetes.UninstallRke2NginxConfig(dryRun)
	}
}

func helmUninstall(kubeconfig string, deployment string, filter string, dryRun bool) string {
	jsonpath := fmt.Sprintf("jsonpath={.items[?(@.metadata.name==\"%s\")].metadata.namespace}", deployment)
	args := []string{"get", "-A", "deploy", "-o", jsonpath}
	if filter != "" {
		args = append(args, filter)
	}

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", args...)
	if err != nil {
		log.Info().Err(err).Msgf("Failed to find %s's namespace, skipping removal", deployment)
	}
	namespace := string(out)
	if namespace != "" {
		helmArgs := []string{}
		if kubeconfig != "" {
			helmArgs = append(helmArgs, "--kubeconfig", kubeconfig)
		}
		helmArgs = append(helmArgs, "uninstall", "-n", namespace, deployment)

		if dryRun {
			log.Info().Msgf("Would run helm %s", strings.Join(helmArgs, " "))
		} else {
			log.Info().Msgf("Uninstalling %s", deployment)
			message := "Failed to run helm " + strings.Join(helmArgs, " ")
			err := utils.RunCmd("helm", helmArgs...)
			if err != nil {
				log.Fatal().Err(err).Msg(message)
			}
		}
	}
	return namespace
}

const kubernetesHelp = `
Note that removing the volumes could also be handled automatically depending on the StorageClass used
when installed on a kubernetes cluster.

For instance on a default K3S install, the local-path-provider storage volumes will
be automatically removed when deleting the deployment even if --purge-volumes argument is not used.`
