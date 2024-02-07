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

type setConfigChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Keys          []string
	ConfigChannelLabels          []string
}

func setConfigChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setConfigChannels",
		Short: "Replace the existing set of
 configuration channels on the given activation keys.
 Channels are ranked by their order in the array.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setConfigChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setConfigChannels)
		},
	}

	cmd.Flags().String("Keys", "", "$desc")
	cmd.Flags().String("ConfigChannelLabels", "", "$desc")

	return cmd
}

func setConfigChannels(globalFlags *types.GlobalFlags, flags *setConfigChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Keys, flags.ConfigChannelLabels)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

