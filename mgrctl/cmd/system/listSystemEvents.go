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

type listSystemEventsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	ActionType          string
	EarliestDate          $date
}

func listSystemEventsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSystemEvents",
		Short: "List system actions of the specified type that were *scheduled* against the given server after the
 specified date. "actionType" should be exactly the string returned in the action_type field
 from the listSystemEvents(sessionKey, serverId) method. For example,
 'Package Install' or 'Initiate a kickstart for a virtual guest.'
 Note: see also system.getEventHistory method which returns a history of all events.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSystemEventsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSystemEvents)
		},
	}

	cmd.Flags().String("Sid", "", "ID of system.")
	cmd.Flags().String("ActionType", "", "Type of the action.")
	cmd.Flags().String("EarliestDate", "", "")

	return cmd
}

func listSystemEvents(globalFlags *types.GlobalFlags, flags *listSystemEventsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.ActionType, flags.EarliestDate)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

