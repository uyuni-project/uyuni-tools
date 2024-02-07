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

type getVersionFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func getVersionCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getVersion",
		Short: "Returns the version of the API.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getVersionFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getVersion)
		},
	}


	return cmd
}

func getVersion(globalFlags *types.GlobalFlags, flags *getVersionFlags, cmd *cobra.Command, args []string) error {

res, err := api.Api(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

