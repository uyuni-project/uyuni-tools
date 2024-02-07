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

type listSubscribedChildChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func listSubscribedChildChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSubscribedChildChannels",
		Short: "Returns a list of subscribed child channels.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSubscribedChildChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSubscribedChildChannels)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listSubscribedChildChannels(globalFlags *types.GlobalFlags, flags *listSubscribedChildChannelsFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
