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

type listMyChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listMyChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listMyChannels",
		Short: "List all software channels that belong to the user's organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listMyChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listMyChannels)
		},
	}


	return cmd
}

func listMyChannels(globalFlags *types.GlobalFlags, flags *listMyChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := channel.Channel(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

