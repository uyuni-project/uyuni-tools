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

type hostnameFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SearchTerm          string
}

func hostnameCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hostname",
		Short: "List the systems which match this hostname",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags hostnameFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, hostname)
		},
	}

	cmd.Flags().String("SearchTerm", "", "")

	return cmd
}

func hostname(globalFlags *types.GlobalFlags, flags *hostnameFlags, cmd *cobra.Command, args []string) error {

res, err := search.Search(&flags.ConnectionDetails, flags.SearchTerm)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

