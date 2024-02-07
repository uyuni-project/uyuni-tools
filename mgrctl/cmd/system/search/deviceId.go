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

type deviceIdFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SearchTerm          string
}

func deviceIdCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deviceId",
		Short: "List the systems which match this device id",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deviceIdFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deviceId)
		},
	}

	cmd.Flags().String("SearchTerm", "", "")

	return cmd
}

func deviceId(globalFlags *types.GlobalFlags, flags *deviceIdFlags, cmd *cobra.Command, args []string) error {

res, err := search.Search(&flags.ConnectionDetails, flags.SearchTerm)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

