package api

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/api"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getApiCallListFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func getApiCallListCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getApiCallList",
		Short: "Lists all available api calls grouped by namespace",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getApiCallListFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getApiCallList)
		},
	}

	return cmd
}

func getApiCallList(globalFlags *types.GlobalFlags, flags *getApiCallListFlags, cmd *cobra.Command, args []string) error {

	res, err := api.Api(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
