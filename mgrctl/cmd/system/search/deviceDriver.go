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

type deviceDriverFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SearchTerm          string
}

func deviceDriverCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deviceDriver",
		Short: "List the systems which match this device driver.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deviceDriverFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deviceDriver)
		},
	}

	cmd.Flags().String("SearchTerm", "", "")

	return cmd
}

func deviceDriver(globalFlags *types.GlobalFlags, flags *deviceDriverFlags, cmd *cobra.Command, args []string) error {

res, err := search.Search(&flags.ConnectionDetails, flags.SearchTerm)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

