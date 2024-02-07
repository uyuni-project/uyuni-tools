package search

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/packages/search"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type nameAndSummaryFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Query          string
}

func nameAndSummaryCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nameAndSummary",
		Short: "Search the lucene package indexes for all packages which
          match the given query in name or summary.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags nameAndSummaryFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nameAndSummary)
		},
	}

	cmd.Flags().String("Query", "", "text to match in package name or summary")

	return cmd
}

func nameAndSummary(globalFlags *types.GlobalFlags, flags *nameAndSummaryFlags, cmd *cobra.Command, args []string) error {

res, err := search.Search(&flags.ConnectionDetails, flags.Query)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

