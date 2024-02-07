package content

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/sync/content"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type addChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	MirrorUrl          string
}

func addChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addChannels",
		Short: "Add a new channel to the #product() database",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addChannels)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "Label of the channel to add")
	cmd.Flags().String("MirrorUrl", "", "Sync from mirror temporarily")

	return cmd
}

func addChannels(globalFlags *types.GlobalFlags, flags *addChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := content.Content(&flags.ConnectionDetails, flags.ChannelLabel, flags.MirrorUrl)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

