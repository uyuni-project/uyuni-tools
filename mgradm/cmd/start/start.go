// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package start

import (
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type startFlags struct {
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[startFlags]) *cobra.Command {
	startCmd := &cobra.Command{
		Use:     "start",
		GroupID: "management",
		Short:   L("Start the server"),
		Long:    L("Start the server"),
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags startFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}
	startCmd.SetUsageTemplate(startCmd.UsageTemplate())

	return startCmd
}

// NewCommand starts the server.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, podmanStart)
}
