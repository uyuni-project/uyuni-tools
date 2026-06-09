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

type getOptions struct {
	api.ConnectionDetails `mapstructure:"api"`
	OutputFormat          string `mapstructure:"output"` // e.g., "json", "yaml", "table"
}

// NewCommand generates the 'mgrctl get' command root.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var flags getOptions

	getCmd := &cobra.Command{
		Use:   "get [resource] [options]",
		Short: L("Display one or many resources (systems, groups, etc.)"),
		Long:  L("Prints tables, JSON, or YAML of Uyuni API resources. Modeled after kubectl get."),
	}

	getCmd.PersistentFlags().StringVarP(
		&flags.OutputFormat,
		"output",
		"o",
		"table",
		L("Output format. One of: table|json|yaml"),
	)

	// passing the -o flag to the subcommands to be able to use it in the same way for all the resources
	getCmd.AddCommand(newSystemCommand(globalFlags, &flags))

	api.AddAPIFlags(getCmd)
	return getCmd
}
