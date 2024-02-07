package systemgroup

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/systemgroup"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type scheduleApplyErrataToActiveFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemGroupName          string
	ErrataIds          []int
	EarliestOccurrence          $date
	OnlyRelevant          bool
}

func scheduleApplyErrataToActiveCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleApplyErrataToActive",
		Short: "Schedules an action to apply errata updates to active systems
 from a group.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleApplyErrataToActiveFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleApplyErrataToActive)
		},
	}

	cmd.Flags().String("SystemGroupName", "", "")
	cmd.Flags().String("ErrataIds", "", "$desc")
	cmd.Flags().String("EarliestOccurrence", "", "")
	cmd.Flags().String("OnlyRelevant", "", "")

	return cmd
}

func scheduleApplyErrataToActive(globalFlags *types.GlobalFlags, flags *scheduleApplyErrataToActiveFlags, cmd *cobra.Command, args []string) error {

res, err := systemgroup.Systemgroup(&flags.ConnectionDetails, flags.SystemGroupName, flags.ErrataIds, flags.EarliestOccurrence, flags.OnlyRelevant)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

