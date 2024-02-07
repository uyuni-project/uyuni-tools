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

type channelExistsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
}

func channelExistsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channelExists",
		Short: "Check for the existence of the config channel provided.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags channelExistsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, channelExists)
		},
	}

	cmd.Flags().String("Label", "", "channel to check for")

	return cmd
}

func channelExists(globalFlags *types.GlobalFlags, flags *channelExistsFlags, cmd *cobra.Command, args []string) error {

res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

