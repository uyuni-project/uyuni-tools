
package maintenance

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "maintenance",
		Short: "Provides methods to access and modify Maintenance Schedules related entities",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(retractScheduleFromSystemsCommand(globalFlags))
	cmd.AddCommand(getScheduleDetailsCommand(globalFlags))
	cmd.AddCommand(updateScheduleCommand(globalFlags))
	cmd.AddCommand(listSystemsWithScheduleCommand(globalFlags))
	cmd.AddCommand(listCalendarLabelsCommand(globalFlags))
	cmd.AddCommand(listScheduleNamesCommand(globalFlags))
	cmd.AddCommand(assignScheduleToSystemsCommand(globalFlags))
	cmd.AddCommand(createCalendarCommand(globalFlags))
	cmd.AddCommand(deleteCalendarCommand(globalFlags))
	cmd.AddCommand(deleteScheduleCommand(globalFlags))
	cmd.AddCommand(createScheduleCommand(globalFlags))
	cmd.AddCommand(getCalendarDetailsCommand(globalFlags))
	cmd.AddCommand(createCalendarWithUrlCommand(globalFlags))
	cmd.AddCommand(updateCalendarCommand(globalFlags))
	cmd.AddCommand(refreshCalendarCommand(globalFlags))

	return cmd
}
