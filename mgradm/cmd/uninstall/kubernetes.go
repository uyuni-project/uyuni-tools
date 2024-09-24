// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package uninstall

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func uninstallForKubernetes(
	globalFlags *types.GlobalFlags,
	flags *utils.UninstallFlags,
	cmd *cobra.Command,
	args []string,
) error {
	if flags.Purge.Volumes {
		log.Warn().Msg(L("--purge-volumes is ignored on a kubernetes deployment"))
	}
	if flags.Purge.Images {
		log.Warn().Msg(L("--purge-images is ignored on a kubernetes deployment"))
	}

	clusterInfos, err := kubernetes.CheckCluster()
	if err != nil {
		return err
	}
	kubeconfig := clusterInfos.GetKubeconfig()

	// TODO Find all the PVs related to the server if we want to delete them

	// Uninstall uyuni
	serverConnection := shared.NewConnection("kubectl", "", kubernetes.ServerFilter)
	serverNamespace, err := serverConnection.GetNamespace("")
	if err != nil {
		return err
	}
	if err := kubernetes.HelmUninstall(serverNamespace, kubeconfig, kubernetes.ServerApp, !flags.Force); err != nil {
		return err
	}

	// Remove the remaining configmap and secrets
	if serverNamespace != "" {
		_, err := utils.RunCmdOutput(zerolog.TraceLevel, "kubectl", "-n", serverNamespace, "get", "secret", "uyuni-ca")
		caSecret := "uyuni-ca"
		if err != nil {
			caSecret = ""
		}

		if !flags.Force {
			log.Info().Msgf(L("Would run %s"), fmt.Sprintf("kubectl delete -n %s configmap uyuni-ca", serverNamespace))
			log.Info().Msgf(L("Would run %s"), fmt.Sprintf("kubectl delete -n %s secret uyuni-cert %s", serverNamespace, caSecret))
		} else {
			log.Info().Msgf(L("Running %s"), fmt.Sprintf("kubectl delete -n %s configmap uyuni-ca", serverNamespace))
			if err := utils.RunCmd("kubectl", "delete", "-n", serverNamespace, "configmap", "uyuni-ca"); err != nil {
				log.Info().Err(err).Msgf(L("Failed deleting config map"))
			}

			log.Info().Msgf(L("Running %s"), fmt.Sprintf("kubectl delete -n %s secret uyuni-cert %s", serverNamespace, caSecret))

			args := []string{"delete", "-n", serverNamespace, "secret", "uyuni-cert"}
			if caSecret != "" {
				args = append(args, caSecret)
			}
			err := utils.RunCmd("kubectl", args...)
			if err != nil {
				log.Info().Err(err).Msgf(L("Failed deleting secret"))
			}
		}
	}

	// TODO Remove the PVs or wait for their automatic removal if purge is requested
	// Also wait if the PVs are dynamic with Delete reclaim policy but the user didn't ask to purge them
	// Since some storage plugins don't handle Delete policy, we may need to check for error events to avoid infinite loop

	// Uninstall cert-manager if we installed it
	certManagerConnection := shared.NewConnection("kubectl", "", "")
	// TODO: re-add "-linstalledby=mgradm" filter once the label is added in helm release
	// mgradm/shared/kubernetes/certificates.go:124 was supposed to be addressing it
	certManagerNamespace, err := certManagerConnection.GetNamespace("cert-manager")
	if err != nil {
		return err
	}
	if err := kubernetes.HelmUninstall(certManagerNamespace, kubeconfig, "cert-manager", !flags.Force); err != nil {
		return err
	}

	// Remove the K3s Traefik config
	if clusterInfos.IsK3s() {
		kubernetes.UninstallK3sTraefikConfig(!flags.Force)
	}

	// Remove the rke2 nginx config
	if clusterInfos.IsRke2() {
		kubernetes.UninstallRke2NginxConfig(!flags.Force)
	}

	if !flags.Force {
		log.Warn().Msg(L("Nothing has been uninstalled, run with --force to actually uninstall"))
	}
	log.Warn().Msg(L("Volumes have not been touched. Depending on the storage class used, they may not have been removed"))
	return nil
}
