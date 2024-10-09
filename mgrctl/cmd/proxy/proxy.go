// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand entry command for managing cache.
// Setup for subcommand to clear (the cache).
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "proxy",
		Short: L("Manage proxy configurations"),
		Long:  L("Manage proxy configurations"),
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: L("Create proxy configurations"),
		Long:  L("Create proxy configurations"),
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	createCmd.AddCommand(NewConfigCommand(globalFlags))

	cmd.AddCommand(createCmd)
	return cmd
}
