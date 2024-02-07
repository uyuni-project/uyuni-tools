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

type removeConfigChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Keys                  []string
	ConfigChannelLabels   []string
}

func removeConfigChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeConfigChannels",
		Short: "Remove configuration channels from the given activation keys.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeConfigChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeConfigChannels)
		},
	}

	cmd.Flags().String("Keys", "", "$desc")
	cmd.Flags().String("ConfigChannelLabels", "", "$desc")

	return cmd
}

func removeConfigChannels(globalFlags *types.GlobalFlags, flags *removeConfigChannelsFlags, cmd *cobra.Command, args []string) error {

	res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Keys, flags.ConfigChannelLabels)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
