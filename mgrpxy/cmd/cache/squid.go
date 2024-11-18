// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewCommand entry command for managing cache.
// Setup for subcommand to clear (the cache).
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var cacheCmd = &cobra.Command{
		Use:   "cache",
		Short: L("Manage proxy cache"),
		Long:  L("Manage proxy cache"),
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
	}

	cacheCmd.AddCommand(NewClearCmd(globalFlags))
	return cacheCmd
}
