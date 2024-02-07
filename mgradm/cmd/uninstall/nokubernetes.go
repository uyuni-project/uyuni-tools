// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build nok8s

package uninstall

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func uninstallForKubernetes(
	globalFlags *types.GlobalFlags,
	flags *uninstallFlags,
	cmd *cobra.Command,
	args []string,
) error {
	return nil
}

const kubernetesHelp = ""
