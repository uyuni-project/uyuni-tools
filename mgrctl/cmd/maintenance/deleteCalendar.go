package maintenance

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/maintenance"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type deleteCalendarFlags struct {
	api.ConnectionDetails  `mapstructure:"api"`
	Label                  string
	CancelScheduledActions bool
}

func deleteCalendarCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteCalendar",
		Short: "Remove a maintenance calendar",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteCalendarFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteCalendar)
		},
	}

	cmd.Flags().String("Label", "", "maintenance calendar label")
	cmd.Flags().String("CancelScheduledActions", "", "cancel actions of affected schedules")

	return cmd
}

func deleteCalendar(globalFlags *types.GlobalFlags, flags *deleteCalendarFlags, cmd *cobra.Command, args []string) error {

	res, err := maintenance.Maintenance(&flags.ConnectionDetails, flags.Label, flags.CancelScheduledActions)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
