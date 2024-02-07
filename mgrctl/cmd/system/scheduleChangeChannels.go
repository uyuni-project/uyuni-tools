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

type scheduleChangeChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	BaseChannelLabel          string
	ChildLabels          []string
	EarliestOccurrence          $type
	Sids          []int
}

func scheduleChangeChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleChangeChannels",
		Short: "Schedule an action to change the channels of the given system. Works for both traditional
 and Salt systems.
 This method accepts labels for the base and child channels.
 If the user provides an empty string for the channelLabel, the current base channel and
 all child channels will be removed from the system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleChangeChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleChangeChannels)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("BaseChannelLabel", "", "")
	cmd.Flags().String("ChildLabels", "", "$desc")
	cmd.Flags().String("EarliestOccurrence", "", "the time/date to schedule the action")
	cmd.Flags().String("Sids", "", "$desc")

	return cmd
}

func scheduleChangeChannels(globalFlags *types.GlobalFlags, flags *scheduleChangeChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.BaseChannelLabel, flags.ChildLabels, flags.EarliestOccurrence, flags.Sids)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

