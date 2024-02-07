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

type listPopularChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	PopularityCount          int
}

func listPopularChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listPopularChannels",
		Short: "List the most popular software channels.  Channels that have at least
 the number of systems subscribed as specified by the popularity count will be
 returned.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listPopularChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listPopularChannels)
		},
	}

	cmd.Flags().String("PopularityCount", "", "")

	return cmd
}

func listPopularChannels(globalFlags *types.GlobalFlags, flags *listPopularChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := channel.Channel(&flags.ConnectionDetails, flags.PopularityCount)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

