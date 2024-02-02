// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package uninstall

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func uninstallForKubernetes(
	globalFlags *types.GlobalFlags,
	flags *uninstallFlags,
	cmd *cobra.Command,
	args []string,
) error {
	clusterInfos := kubernetes.CheckCluster()
	kubeconfig := clusterInfos.GetKubeconfig()

	// TODO Find all the PVs related to the server if we want to delete them

	// Uninstall uyuni
	namespace := kubernetes.HelmUninstall(kubeconfig, "uyuni", "", flags.DryRun)

	// Remove the remaining configmap and secrets
	if namespace != "" {

		_, err := utils.RunCmdOutput(zerolog.TraceLevel, "kubectl", "-n", namespace, "get", "secret", "uyuni-ca")
		caSecret := "uyuni-ca"
		if err != nil {
			caSecret = ""
		}

		if flags.DryRun {
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
	kubernetes.HelmUninstall(kubeconfig, "cert-manager", "-linstalledby=mgradm", flags.DryRun)

	// Remove the K3s Traefik config
	if clusterInfos.IsK3s() {
		kubernetes.UninstallK3sTraefikConfig(flags.DryRun)
	}

	// Remove the rke2 nginx config
	if clusterInfos.IsRke2() {
		kubernetes.UninstallRke2NginxConfig(flags.DryRun)
	}
	return nil
}

const kubernetesHelp = kubernetes.UninstallHelp
