// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build nok8s

package uninstall

const kubernetesBuilt = false

func uninstallForKubernetes(dryRun bool) {
}

func helmUninstall(kubeconfig string, deployment string, filter string, dryRun bool) string {
	return ""
}

const kubernetesHelp = ""
