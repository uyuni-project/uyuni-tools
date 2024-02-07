package channel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listManageableChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listManageableChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listManageableChannels",
		Short: "List all software channels that the user is entitled to manage.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listManageableChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listManageableChannels)
		},
	}


	return cmd
}

func listManageableChannels(globalFlags *types.GlobalFlags, flags *listManageableChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := channel.Channel(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

