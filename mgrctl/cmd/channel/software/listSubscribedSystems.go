package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listSubscribedSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func listSubscribedSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSubscribedSystems",
		Short: "Returns list of subscribed systems for the given channel label",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSubscribedSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSubscribedSystems)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "channel to query")

	return cmd
}

func listSubscribedSystems(globalFlags *types.GlobalFlags, flags *listSubscribedSystemsFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

