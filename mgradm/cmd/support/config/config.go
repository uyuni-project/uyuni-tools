// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type configFlags struct {
	Output  string
	Backend string
}

// NewCommand is the command for creates supportconfig.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: L("Extract configuration and logs"),
		Long: L(`Extract the host or cluster configuration and logs as well as those from 
the containers for support to help debugging.`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags configFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, extract)
		},
	}

	configCmd.Flags().StringP("output", "o", "supportconfig.tar.gz", L("path where to extract the data"))
	utils.AddBackendFlag(configCmd)

	return configCmd
}
