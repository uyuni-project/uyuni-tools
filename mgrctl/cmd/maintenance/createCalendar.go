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

type createCalendarFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	Ical          string
}

func createCalendarCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createCalendar",
		Short: "Create a new maintenance calendar",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createCalendarFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createCalendar)
		},
	}

	cmd.Flags().String("Label", "", "maintenance calendar label")
	cmd.Flags().String("Ical", "", "ICal calendar data")

	return cmd
}

func createCalendar(globalFlags *types.GlobalFlags, flags *createCalendarFlags, cmd *cobra.Command, args []string) error {

res, err := maintenance.Maintenance(&flags.ConnectionDetails, flags.Label, flags.Ical)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

