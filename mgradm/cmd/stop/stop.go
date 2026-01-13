// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package stop

import (
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type stopFlags struct {
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[stopFlags]) *cobra.Command {
	stopCmd := &cobra.Command{
		Use:     "stop",
		GroupID: "management",
		Short:   L("Stop the server"),
		Long:    L("Stop the server"),
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags stopFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	stopCmd.SetUsageTemplate(stopCmd.UsageTemplate())

	return stopCmd
}

// NewCommand to stop server.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, podmanStop)
}
