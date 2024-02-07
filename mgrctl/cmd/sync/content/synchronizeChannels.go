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

type synchronizeChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	MirrorUrl          string
}

func synchronizeChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "synchronizeChannels",
		Short: "(Deprecated) Synchronize channels between the Customer Center
             and the #product() database.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags synchronizeChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, synchronizeChannels)
		},
	}

	cmd.Flags().String("MirrorUrl", "", "Sync from mirror temporarily")

	return cmd
}

func synchronizeChannels(globalFlags *types.GlobalFlags, flags *synchronizeChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := content.Content(&flags.ConnectionDetails, flags.MirrorUrl)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

