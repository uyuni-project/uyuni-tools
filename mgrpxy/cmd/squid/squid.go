// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package squid

import (
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand entry command for managing squid cache.
// Setup for subcommand to clear (the cache).
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var squidCmd = &cobra.Command{
		Use:   "squid",
		Short: L("Manage Squid cache"),
		Long:  L("Manage Squid cache"),
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	squidCmd.AddCommand(NewClearCmd(globalFlags))
	return squidCmd
}
