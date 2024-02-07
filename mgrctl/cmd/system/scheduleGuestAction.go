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

type scheduleGuestActionFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	State          string
	Date          $type
}

func scheduleGuestActionCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleGuestAction",
		Short: "Schedules a guest action for the specified virtual guest for a given
          date/time.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleGuestActionFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleGuestAction)
		},
	}

	cmd.Flags().String("Sid", "", "the system Id of the guest")
	cmd.Flags().String("State", "", "One of the following actions  'start',          'suspend', 'resume', 'restart', 'shutdown'.")
	cmd.Flags().String("Date", "", "the time/date to schedule the action")

	return cmd
}

func scheduleGuestAction(globalFlags *types.GlobalFlags, flags *scheduleGuestActionFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.State, flags.Date)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

