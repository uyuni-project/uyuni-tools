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

type uuidFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SearchTerm          string
}

func uuidCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uuid",
		Short: "List the systems which match this UUID",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags uuidFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, uuid)
		},
	}

	cmd.Flags().String("SearchTerm", "", "")

	return cmd
}

func uuid(globalFlags *types.GlobalFlags, flags *uuidFlags, cmd *cobra.Command, args []string) error {

res, err := search.Search(&flags.ConnectionDetails, flags.SearchTerm)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

