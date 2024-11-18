// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
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
	dryRun := !flags.Force

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
	cnx := shared.NewConnection("kubectl", "", kubernetes.ProxyFilter)
	namespace, err := cnx.GetNamespace("")
	if err != nil {
		return err
	}
	if err := kubernetes.HelmUninstall(namespace, kubeconfig, kubernetes.ProxyApp, dryRun); err != nil {
		return err
	}

	// TODO Remove the PVs or wait for their automatic removal if purge is requested
	// Also wait if the PVs are dynamic with Delete reclaim policy but the user didn't ask to purge them
	// Since some storage plugins don't handle Delete policy, we may need to check for error events to avoid infinite loop

	// Remove the K3s Traefik config
	if clusterInfos.IsK3s() {
		kubernetes.UninstallK3sTraefikConfig(dryRun)
	}

	// Remove the rke2 nginx config
	if clusterInfos.IsRke2() {
		kubernetes.UninstallRke2NginxConfig(dryRun)
	}

	if dryRun {
		log.Warn().Msg(L("Nothing has been uninstalled, run with --force to actually uninstall"))
	}
	log.Warn().Msg(L("Volumes have not been touched. Depending on the storage class used, they may not have been removed"))
	return nil
}
