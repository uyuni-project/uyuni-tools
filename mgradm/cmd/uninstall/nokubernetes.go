// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build nok8s

package uninstall

func uninstallForKubernetes(
	globalFlags *types.GlobalFlags,
	flags *uninstallFlags,
	cmd *cobra.Command,
	args []string,
) {
}

const kubernetesHelp = ""
