package provider

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/packages/provider"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all Package Providers.
 User executing the request must be a #product() administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, list)
		},
	}


	return cmd
}

func list(globalFlags *types.GlobalFlags, flags *listFlags, cmd *cobra.Command, args []string) error {

res, err := provider.Provider(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

