package activationkey

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/activationkey"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type removeChildChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Key                   string
	ChildChannelLabels    []string
}

func removeChildChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeChildChannels",
		Short: "Remove child channels from an activation key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeChildChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeChildChannels)
		},
	}

	cmd.Flags().String("Key", "", "")
	cmd.Flags().String("ChildChannelLabels", "", "$desc")

	return cmd
}

func removeChildChannels(globalFlags *types.GlobalFlags, flags *removeChildChannelsFlags, cmd *cobra.Command, args []string) error {

	res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Key, flags.ChildChannelLabels)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
