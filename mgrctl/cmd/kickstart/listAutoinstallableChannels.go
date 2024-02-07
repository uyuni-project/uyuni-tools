package kickstart

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listAutoinstallableChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listAutoinstallableChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAutoinstallableChannels",
		Short: "List autoinstallable channels for the logged in user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAutoinstallableChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAutoinstallableChannels)
		},
	}


	return cmd
}

func listAutoinstallableChannels(globalFlags *types.GlobalFlags, flags *listAutoinstallableChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := kickstart.Kickstart(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

