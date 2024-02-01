// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build nok8s

package uninstall

func uninstallForKubernetes(dryRun bool) {
}

func helmUninstall(kubeconfig string, deployment string, filter string, dryRun bool) string {
	return ""
}

const kubernetesHelp = ""
