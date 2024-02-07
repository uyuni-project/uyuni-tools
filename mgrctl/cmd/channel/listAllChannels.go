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

type listAllChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listAllChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAllChannels",
		Short: "List all software channels that the user's organization is entitled to.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAllChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAllChannels)
		},
	}


	return cmd
}

func listAllChannels(globalFlags *types.GlobalFlags, flags *listAllChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := channel.Channel(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

