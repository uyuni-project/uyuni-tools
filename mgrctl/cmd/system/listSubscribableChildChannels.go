package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listSubscribableChildChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func listSubscribableChildChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSubscribableChildChannels",
		Short: "Returns a list of subscribable child channels.  This only shows channels
 the system is *not* currently subscribed to.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSubscribableChildChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSubscribableChildChannels)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listSubscribableChildChannels(globalFlags *types.GlobalFlags, flags *listSubscribableChildChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

