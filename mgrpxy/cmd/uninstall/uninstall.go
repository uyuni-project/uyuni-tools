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

// NewCommand for uninstall proxy.
func NewCommand(globalFlags *types.GlobalFlags) (*cobra.Command, error) {
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: L("Uninstall a proxy"),
		Long: L(`Uninstall a proxy and optionally the corresponding volumes.
By default it will only print what would be done, use --force to actually remove.`) + kubernetes.UninstallHelp(),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags utils.UninstallFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, uninstall)
		},
	}
	utils.AddUninstallFlags(uninstallCmd, true)

	return uninstallCmd, nil
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
