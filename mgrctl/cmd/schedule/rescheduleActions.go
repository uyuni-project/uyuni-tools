package schedule

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/schedule"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type rescheduleActionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ActionIds          []int
	OnlyFailed          bool
}

func rescheduleActionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rescheduleActions",
		Short: "Reschedule all actions in the given list.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags rescheduleActionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, rescheduleActions)
		},
	}

	cmd.Flags().String("ActionIds", "", "$desc")
	cmd.Flags().String("OnlyFailed", "", "True to only reschedule failed actions, False to reschedule all")

	return cmd
}

func rescheduleActions(globalFlags *types.GlobalFlags, flags *rescheduleActionsFlags, cmd *cobra.Command, args []string) error {

res, err := schedule.Schedule(&flags.ConnectionDetails, flags.ActionIds, flags.OnlyFailed)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

