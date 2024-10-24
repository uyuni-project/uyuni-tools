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

type sqlFlags struct {
	Database       string
	Interactive    bool
	ForceOverwrite bool   `mapstructure:"force"`
	OutputFile     string `mapstructure:"output"`
	Backend        string
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[sqlFlags]) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sql [sql-file]",
		Short: L("Execute SQL query"),
		Long: L(`Execute SQL query either provided in sql-file or passed through standard input.

Examples:

  Run the 'select hostname from rhnserver;' query using echo:

  # echo 'select hostname from rhnserver;' | mgradm support sql

  Run in interative mode:

  # mgradm support sql -i

  Running the SQL queries in example.sql file and output them to out.log file

  # mgradm support sql example.sql -o out.log

`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags sqlFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	cmd.Flags().StringP("database", "d", "productdb", L("Target database, can be 'reportdb' or 'productdb'"))
	cmd.Flags().BoolP("interactive", "i", false, L("Start in interactive mode"))
	cmd.Flags().BoolP("force", "f", false, L("Force overwrite of output file if already exists"))
	cmd.Flags().StringP("output", "o", "", L("Write output to the file instead of standard output"))
	utils.AddBackendFlag(cmd)

	return cmd
}

// NewCommand adds support sql command.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, doSql)
}
