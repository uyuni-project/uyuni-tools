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

type deviceDescriptionFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SearchTerm          string
}

func deviceDescriptionCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deviceDescription",
		Short: "List the systems which match the device description.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deviceDescriptionFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deviceDescription)
		},
	}

	cmd.Flags().String("SearchTerm", "", "")

	return cmd
}

func deviceDescription(globalFlags *types.GlobalFlags, flags *deviceDescriptionFlags, cmd *cobra.Command, args []string) error {

res, err := search.Search(&flags.ConnectionDetails, flags.SearchTerm)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

