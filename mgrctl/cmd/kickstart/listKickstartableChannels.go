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

type listKickstartableChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listKickstartableChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listKickstartableChannels",
		Short: "List kickstartable channels for the logged in user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listKickstartableChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listKickstartableChannels)
		},
	}


	return cmd
}

func listKickstartableChannels(globalFlags *types.GlobalFlags, flags *listKickstartableChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := kickstart.Kickstart(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

