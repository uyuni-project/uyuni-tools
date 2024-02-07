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

type listSystemChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func listSystemChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSystemChannels",
		Short: "Returns a list of channels that a system is subscribed to for the
 given system id",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSystemChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSystemChannels)
		},
	}

	cmd.Flags().String("Sid", "", "system ID")

	return cmd
}

func listSystemChannels(globalFlags *types.GlobalFlags, flags *listSystemChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

