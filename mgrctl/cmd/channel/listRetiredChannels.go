package channel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listRetiredChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listRetiredChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listRetiredChannels",
		Short: "List all retired software channels.  These are channels that the user's
 organization is entitled to, but are no longer supported because they have reached
 their 'end-of-life' date.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listRetiredChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listRetiredChannels)
		},
	}


	return cmd
}

func listRetiredChannels(globalFlags *types.GlobalFlags, flags *listRetiredChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := channel.Channel(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

