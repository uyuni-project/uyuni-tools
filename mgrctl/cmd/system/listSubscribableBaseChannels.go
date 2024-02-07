package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listSubscribableBaseChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func listSubscribableBaseChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSubscribableBaseChannels",
		Short: "Returns a list of subscribable base channels.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSubscribableBaseChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSubscribableBaseChannels)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listSubscribableBaseChannels(globalFlags *types.GlobalFlags, flags *listSubscribableBaseChannelsFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
