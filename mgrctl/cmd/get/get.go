// SPDX-FileCopyrightText: 2026 Jayprakash
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OutputFormat          string `mapstructure:"output"`
	Filter                string `mapstructure:"filter"`
	Page                  int    `mapstructure:"-"`
	PageSize              int    `mapstructure:"-"`
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var flags getFlags

	cmd := &cobra.Command{
		Use:   "get <resource-type> [flags]",
		Short: L("Display one or many resources"),
		Example: `  mgrctl get system
  mgrctl get systemgroup
  mgrctl get system --page 0 --page-size 10
  mgrctl get system --filter "extra_pkg_count>0"`,
		Long: L(`Fetch and display Uyuni API resources in table, JSON, or YAML format.

The resource type is passed as the first argument. Filtering and pagination
are handled server-side when the API endpoint supports it.

Available resource types:
` + GetResourceHelp() + `

Filter Operators:
  >=, <=, !=, =, >, < (e.g. extra_pkg_count>0)

Custom Columns:
  Use -o custom-columns=HEADER:path,HEADER:path to extract specific JSON fields.
  Example: -o custom-columns=ID:.id,NAME:.name`),
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs: registeredTypes(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, runGet)
		},
	}

	utils.AddOutputFlag(cmd, &flags.OutputFormat)
	cmd.Flags().StringVar(&flags.Filter, "filter", "",
		L("Filter expression (e.g. extra_pkg_count=0)."+
			" See https://www.uyuni-project.org/uyuni-docs-api/uyuni/index.html for available keys"))
	cmd.Flags().IntVar(&flags.Page, "page", 0,
		L("Page number for paginated results (0-indexed). Not all resources support pagination."))
	cmd.Flags().IntVar(&flags.PageSize, "page-size", 50,
		L("Number of items per page. Not all resources support pagination."))

	api.AddAPIFlags(cmd)
	return cmd
}

func runGet(_ *types.GlobalFlags, flags *getFlags, _ *cobra.Command, args []string) error {
	resourceType := args[0]
	log.Debug().Msgf("Running get %s", resourceType)

	resource, err := lookupResource(resourceType)
	if err != nil {
		return err
	}
	client, err := api.Init(&flags.ConnectionDetails)
	if err == nil && (client.Details.User != "" || client.Details.InSession) {
		err = client.Login()
	}
	if err != nil {
		return utils.Errorf(err, L("unable to login to the server"))
	}

	return resource.ListAndPrint(client, flags.Filter, flags.Page, flags.PageSize, flags.OutputFormat, os.Stdout)
}
