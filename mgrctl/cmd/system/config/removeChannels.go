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

type removeChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func removeChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeChannels",
		Short: "Remove config channels from the given servers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeChannels)
		},
	}


	return cmd
}

func removeChannels(globalFlags *types.GlobalFlags, flags *removeChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := config.Config(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

