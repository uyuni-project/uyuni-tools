package content

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/sync/content"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type synchronizeSubscriptionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func synchronizeSubscriptionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "synchronizeSubscriptions",
		Short: "Synchronize subscriptions between the Customer Center
             and the #product() database.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags synchronizeSubscriptionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, synchronizeSubscriptions)
		},
	}


	return cmd
}

func synchronizeSubscriptions(globalFlags *types.GlobalFlags, flags *synchronizeSubscriptionsFlags, cmd *cobra.Command, args []string) error {

res, err := content.Content(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

