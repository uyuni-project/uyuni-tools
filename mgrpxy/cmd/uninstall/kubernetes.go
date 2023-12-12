// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
)

const kubernetesBuilt = true

func uninstallForKubernetes(dryRun bool) {
	clusterInfos := kubernetes.CheckCluster()
	kubeconfig := clusterInfos.GetKubeconfig()

	// TODO Find all the PVs related to the server if we want to delete them

	// Uninstall uyuni
	kubernetes.HelmUninstall(kubeconfig, "uyuni-proxy", "", dryRun)

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
}
