package search

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/search"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type ipFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SearchTerm          string
}

func ipCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ip",
		Short: "List the systems which match this ip.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags ipFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, ip)
		},
	}

	cmd.Flags().String("SearchTerm", "", "")

	return cmd
}

func ip(globalFlags *types.GlobalFlags, flags *ipFlags, cmd *cobra.Command, args []string) error {

res, err := search.Search(&flags.ConnectionDetails, flags.SearchTerm)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

