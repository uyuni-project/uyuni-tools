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

type listSoftwareChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listSoftwareChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSoftwareChannels",
		Short: "List all visible software channels.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSoftwareChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSoftwareChannels)
		},
	}


	return cmd
}

func listSoftwareChannels(globalFlags *types.GlobalFlags, flags *listSoftwareChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := channel.Channel(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

