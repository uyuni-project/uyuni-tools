// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
    "github.com/spf13/cobra"
    "github.com/uyuni-project/uyuni-tools/shared/api"
    . "github.com/uyuni-project/uyuni-tools/shared/l10n"
    "github.com/uyuni-project/uyuni-tools/shared/types"
)

type getFlags struct {
    api.ConnectionDetails `mapstructure:"api"`
    Output                string `mapstructure:"output"`
}

// NewCommand returns the get command.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
    getCmd := &cobra.Command{
        Use:   "get",
        Short: L("Query Uyuni server resources"),
        Long: L(`Query resources from the Uyuni server API.

Supports multiple resource types and output formats.

Examples:
  # List all systems in table format
  mgrctl get system

  # Search for a system by name
  mgrctl get system webserver

  # Output as JSON
  mgrctl get system -o json`),
    }

    getCmd.AddCommand(newSystemCommand(globalFlags))
    api.AddAPIFlags(getCmd)

    return getCmd
}