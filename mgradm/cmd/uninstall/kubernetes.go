// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package uninstall

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func uninstallForKubernetes(
	_ *types.GlobalFlags,
	flags *utils.UninstallFlags,
	_ *cobra.Command,
	_ []string,
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

	// Remove all Uyuni resources
	if serverNamespace != "" {
		objects := "job,deploy,svc,ingress,pvc,cm,secret"
		if kubernetes.HasResource("ingressroutetcps") {
			objects += ",middlewares,ingressroutetcps,ingressrouteudps"
		}

		if kubernetes.HasResource("issuers") {
			objects += ",issuers,certificates"
		}
		deleteCmd := []string{
			"kubectl", "delete", "-n", serverNamespace, objects,
			"-l", kubernetes.AppLabel + "=" + kubernetes.ServerApp,
		}
		if !flags.Force {
			log.Info().Msgf(L("Would run %s"), strings.Join(deleteCmd, " "))
		} else {
			if err := utils.RunCmd(deleteCmd[0], deleteCmd[1:]...); err != nil {
				return utils.Errorf(err, L("failed to delete server resources"))
			}
		}
	}

	// TODO Remove the PVs or wait for their automatic removal if purge is requested
	// Also wait if the PVs are dynamic with Delete reclaim policy but the user didn't ask to purge them
	// Since some storage plugins don't handle Delete policy, we may need to check for error events to avoid infinite loop

	// Uninstall cert-manager if we installed it
	certManagerConnection := shared.NewConnection("kubectl", "", "-linstalledby=mgradm")
	// TODO: re-add "-linstalledby=mgradm" filter once the label is added in helm release
	// mgradm/shared/kubernetes/certificates.go:124 was supposed to be addressing it
	certManagerNamespace, err := certManagerConnection.GetNamespace("cert-manager")
	if err != nil {
		return err
	}
	if certManagerNamespace != "" {
		if err := kubernetes.HelmUninstall(certManagerNamespace, kubeconfig, "cert-manager", !flags.Force); err != nil {
			return err
		}
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
