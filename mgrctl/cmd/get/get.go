// SPDX-FileCopyrightText: 2026 Jay Prakash katara <katarajayprakash@icloud.com>
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// pageFlags groups the pagination flags so they map to the nested
// page.number and page.size configuration keys.
type pageFlags struct {
	Number int `mapstructure:"number"`
	Size   int `mapstructure:"size"`
}

type sortFlags struct {
	By string `mapstructure:"by"`
}

type getFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OutputFormat          string    `mapstructure:"output"`
	Filter                string    `mapstructure:"filter"`
	Page                  pageFlags `mapstructure:"page"`
	Sort                  sortFlags `mapstructure:"sort"`
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var flags getFlags

	cmd := &cobra.Command{
		Use:   "get <resource-type> [flags]",
		Short: L("Display one or many resources"),
		Long: fmt.Sprintf(L(`Display one or many Uyuni API resources.

Prints a table by default. Use --output to get the complete objects as JSON or
YAML, or to select individual fields with custom columns, JSONPath or a Go
template.

Supported Resources:
%[1]s

Filtering:
  Filtering is supported by resources whose API provides it. Currently,
  system supports server-side filtering; systemgroup does not.

  --filter takes a KEY OPERATOR VALUE expression, where OPERATOR is one of
  =, !=, >, <, >=, <=. The available keys depend on the resource type and are
  documented at https://www.uyuni-project.org/uyuni-docs-api/uyuni/index.html

Field names:
  custom-columns uses the Go field names, such as Name or Status.State, while
  jsonpath and go-template use the JSON field names, such as name or
  status.state.

Resource details:
%[2]s`), getResourceHelpSummaries(), getResourceHelpDetails()),
		Example: fmt.Sprintf(L(`  # Print complete objects as JSON or YAML
  mgrctl get system -o json
  mgrctl get system -o yaml

  # Show only the selected columns, using the Go field names
  mgrctl get system -o custom-columns=ID:ID,NAME:Name

  # Read the same column definition from a file, one HEADER:FIELD pair per line
  mgrctl get system -o custom-columns-file=./columns.txt

  # Print a single field of every system, using the JSON field names
  mgrctl get system -o jsonpath='{.items[*].id}'

  # Print one line per system with a Go template
  mgrctl get system -o go-template='{{range .items}}{{.id}} {{.name}}{{"\n"}}{{end}}'

  # Read the same Go template from a file
  mgrctl get system -o go-template-file=./template.tmpl

%s`), getResourceHelpExamples()),
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs: registeredTypes(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, runGet)
		},
	}

	utils.AddOutputFlag(cmd, &flags.OutputFormat)
	cmd.Flags().StringVar(&flags.Filter, "filter", "",
		L("Filter expression, e.g. extra_pkg_count>0. See the operators and keys described above"))
	cmd.Flags().StringVar(&flags.Sort.By, "sort-by", "",
		L("Field to sort results by; supported fields depend on the resource type"))
	cmd.Flags().IntVar(&flags.Page.Number, "page", 0,
		L("Page number for paginated results (0-indexed). "+
			"May be ignored, depending on the resource type and its API"))
	cmd.Flags().IntVar(&flags.Page.Size, "page-size", 50,
		L("Number of items per page. "+
			"May be ignored, depending on the resource type and its API"))

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
	if err != nil {
		return utils.Errorf(err, L("unable to initialize connection to the server"))
	}
	if client.Details.User != "" || client.Details.InSession {
		if err := client.Login(); err != nil {
			return utils.Errorf(err, L("unable to login to the server"))
		}
	}

	return resource.ListAndPrint(client, ListOptions{
		Filter:       flags.Filter,
		SortBy:       flags.Sort.By,
		PageNumber:   flags.Page.Number,
		PageSize:     flags.Page.Size,
		OutputFormat: flags.OutputFormat,
		Out:          os.Stdout,
	})
}
