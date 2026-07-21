// SPDX-FileCopyrightText: 2026 Jayprakash
//
// SPDX-License-Identifier: Apache-2.0

package api_resources

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/get"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flags struct{}

// NewCommand returns a new cobra.Command for listing available API resources.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var f flags
	cmd := &cobra.Command{
		Use:   "api-resources",
		Short: L("Display the supported API resources"),
		Long: L(`Print a table of the API resources supported by 'mgrctl get'.

Each row shows the resource name, its short aliases, and a brief description.`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &f, nil, runAPIResources)
		},
	}
	return cmd
}

func runAPIResources(_ *types.GlobalFlags, _ *flags, _ *cobra.Command, _ []string) error {
	resources := get.GetRegisteredResources()
	if len(resources) == 0 {
		fmt.Println(L("No resources registered."))
		return nil
	}

	fmt.Printf("%-20s %-15s %s\n", "NAME", "ALIASES", "DESCRIPTION")
	for _, res := range resources {
		aliases := ""
		if len(res.Aliases) > 0 {
			aliases = strings.Join(res.Aliases, ", ")
		}
		fmt.Printf("%-20s %-15s %s\n", res.Name, aliases, res.Description)
	}
	return nil
}
