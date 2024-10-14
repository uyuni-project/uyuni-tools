// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[utils.UninstallFlags]) *cobra.Command {
	uninstallCmd := &cobra.Command{
		Use:     "uninstall",
		GroupID: "deploy",
		Short:   L("Uninstall a server"),
		Long: L(`Uninstall a server and optionally the corresponding volumes.
By default it will only print what would be done, use --force to actually remove.`) + kubernetes.UninstallHelp(),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags utils.UninstallFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}
	utils.AddUninstallFlags(uninstallCmd, utils.KubernetesBuilt)

	return uninstallCmd
}

// NewCommand uninstall a server and optionally the corresponding volumes.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, uninstall)
}

func uninstall(
	globalFlags *types.GlobalFlags,
	flags *utils.UninstallFlags,
	cmd *cobra.Command,
	args []string,
) error {
	fn, err := shared.ChoosePodmanOrKubernetes(cmd.Flags(), uninstallForPodman, uninstallForKubernetes)
	if err != nil {
		return err
	}

	return fn(globalFlags, flags, cmd, args)
}
