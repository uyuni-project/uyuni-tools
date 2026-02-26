// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[utils.UninstallFlags]) *cobra.Command {
	uninstallCmd := &cobra.Command{
		Use:     "uninstall",
		GroupID: "deploy",
		Short:   L("Uninstall a proxy"),
		Long: L(`Uninstall a proxy and optionally the corresponding volumes.
By default it will only print what would be done, use --force to actually remove.`),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags utils.UninstallFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}
	utils.AddUninstallFlags(uninstallCmd)

	return uninstallCmd
}

// NewCommand for uninstall proxy.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, uninstallForPodman)
}
