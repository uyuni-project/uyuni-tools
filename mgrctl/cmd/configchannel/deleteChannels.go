package configchannel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/configchannel"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type deleteChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	$param.getFlagName()          $param.getType()
}

func deleteChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteChannels",
		Short: "Delete a list of global config channels.
 Caller must be a config admin.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteChannels)
		},
	}

	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func deleteChannels(globalFlags *types.GlobalFlags, flags *deleteChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

