// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package restart

import (
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type restartFlags struct {
	Backend string
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[restartFlags]) *cobra.Command {
	restartCmd := &cobra.Command{
		Use:     "restart",
		GroupID: "management",
		Short:   L("Restart the server"),
		Long:    L("Restart the server"),
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags restartFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}
	restartCmd.SetUsageTemplate(restartCmd.UsageTemplate())

	return restartCmd
}

// NewCommand to restart server.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, podmanRestart)
}
