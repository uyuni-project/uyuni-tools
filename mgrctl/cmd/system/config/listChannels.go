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

type listChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func listChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listChannels",
		Short: "List all global('Normal', 'State') configuration channels associated to a
              system in the order of their ranking.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listChannels)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listChannels(globalFlags *types.GlobalFlags, flags *listChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := config.Config(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

