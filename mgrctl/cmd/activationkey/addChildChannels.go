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

type addChildChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Key                   string
	ChildChannelLabels    []string
}

func addChildChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addChildChannels",
		Short: "Add child channels to an activation key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addChildChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addChildChannels)
		},
	}

	cmd.Flags().String("Key", "", "")
	cmd.Flags().String("ChildChannelLabels", "", "$desc")

	return cmd
}

func addChildChannels(globalFlags *types.GlobalFlags, flags *addChildChannelsFlags, cmd *cobra.Command, args []string) error {

	res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Key, flags.ChildChannelLabels)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
