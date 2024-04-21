// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type configFlags struct {
	Database       string
	Interactive    bool
	ForceOverwrite bool   `mapstructure:"force"`
	OutputFile     string `mapstructure:"output"`
	Backend        string
}

// Add support sql command.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "sql [sql-file]",
		Short: L("Execute SQL query"),
		Long:  L(`Execute SQL query either provided in sql-file or passed through standard input`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags configFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, doSql)
		},
	}

	configCmd.Flags().StringP("database", "d", "productdb", L("Target database, can be 'reportdb' or 'productdb'"))
	configCmd.Flags().BoolP("interactive", "i", false, L("Start in interactive mode"))
	configCmd.Flags().BoolP("force", "f", false, L("Force overwrite of output file if already exists"))
	configCmd.Flags().StringP("output", "o", "", L("Write output to the file instead of standard output"))
	utils.AddBackendFlag(configCmd)

	return configCmd
}
