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

type nameAndDescriptionFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SearchTerm          string
}

func nameAndDescriptionCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nameAndDescription",
		Short: "List the systems which match this name or description",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags nameAndDescriptionFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nameAndDescription)
		},
	}

	cmd.Flags().String("SearchTerm", "", "")

	return cmd
}

func nameAndDescription(globalFlags *types.GlobalFlags, flags *nameAndDescriptionFlags, cmd *cobra.Command, args []string) error {

res, err := search.Search(&flags.ConnectionDetails, flags.SearchTerm)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

