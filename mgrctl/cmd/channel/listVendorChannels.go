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

type listVendorChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listVendorChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listVendorChannels",
		Short: "Lists all the vendor software channels that the user's organization
 is entitled to.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listVendorChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listVendorChannels)
		},
	}


	return cmd
}

func listVendorChannels(globalFlags *types.GlobalFlags, flags *listVendorChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := channel.Channel(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

