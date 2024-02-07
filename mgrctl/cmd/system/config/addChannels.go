package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/config"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type addChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func addChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addChannels",
		Short: "Given a list of servers and configuration channels,
 this method appends the configuration channels to either the top or
 the bottom (whichever you specify) of a system's subscribed
 configuration channels list. The ordering of the configuration channels
 provided in the add list is maintained while adding.
 If one of the configuration channels in the 'add' list
 has been previously subscribed by a server, the
 subscribed channel will be re-ranked to the appropriate place.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addChannels)
		},
	}


	return cmd
}

func addChannels(globalFlags *types.GlobalFlags, flags *addChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := config.Config(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

