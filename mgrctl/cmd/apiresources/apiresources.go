// SPDX-FileCopyrightText: 2026 Jay Prakash katara <katarajayprakash@icloud.com>
//
// SPDX-License-Identifier: Apache-2.0

package apiresources

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/get"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flags struct{}

type resourceRow struct {
	Name        string
	Aliases     string
	Description string
}

var resourceColumns = []utils.ColumnDef{
	{Header: L("NAME"), Field: "Name"},
	{Header: L("ALIAS NAMES"), Field: "Aliases"},
	{Header: L("DESCRIPTION"), Field: "Description"},
}

// NewCommand returns a new cobra.Command for listing available API resources.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var f flags
	cmd := &cobra.Command{
		Use:   "api-resources",
		Short: L("Display the supported API resources"),
		Long:  L("Print the resource types that can be used with 'mgrctl get'."),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &f, nil, runAPIResources)
		},
	}
	return cmd
}

func runAPIResources(_ *types.GlobalFlags, _ *flags, cmd *cobra.Command, _ []string) error {
	return printResources(cmd.OutOrStdout(), get.GetRegisteredResources())
}

func printResources(out io.Writer, resources []get.ResourceInfo) error {
	if len(resources) == 0 {
		_, err := fmt.Fprintln(out, L("No resources registered."))
		return err
	}

	rows := make([]resourceRow, len(resources))
	for i, resource := range resources {
		rows[i] = resourceRow{
			Name:        resource.Name,
			Aliases:     strings.Join(resource.Aliases, ", "),
			Description: resource.Help.Summary,
		}
	}

	return utils.PrintOutput("table", rows, resourceColumns, out)
}
