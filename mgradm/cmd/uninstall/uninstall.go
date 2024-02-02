// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type uninstallFlags struct {
	Backend      string
	DryRun       bool
	PurgeVolumes bool
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "uninstall a server",
		Long:  "Uninstall a server and optionally the corresponding volumes." + kubernetesHelp,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags uninstallFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, uninstall)
		},
	}
	uninstallCmd.Flags().BoolP("dryRun", "n", false, "Only show what would be done")
	uninstallCmd.Flags().Bool("purgeVolumes", false, "Also remove the volume")

	if utils.KubernetesBuilt {
		utils.AddBackendFlag(uninstallCmd)
	}

	return uninstallCmd
}

func uninstall(
	globalFlags *types.GlobalFlags,
	flags *uninstallFlags,
	cmd *cobra.Command,
	args []string,
) error {
	fn, err := shared.ChoosePodmanOrKubernetes(cmd.Flags(), uninstallForPodman, uninstallForKubernetes)
	if err != nil {
		return err
	}

	return fn(globalFlags, flags, cmd, args)
}
