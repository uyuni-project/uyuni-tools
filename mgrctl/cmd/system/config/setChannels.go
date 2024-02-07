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

type setChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func setChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setChannels",
		Short: "Replace the existing set of config channels on the given servers.
 Channels are ranked according to their order in the configChannelLabels
 array.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setChannels)
		},
	}


	return cmd
}

func setChannels(globalFlags *types.GlobalFlags, flags *setChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := config.Config(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

