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

type scheduleApplyErrataFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sids          []int
	ErrataIds          []int
	AllowModules          bool
	EarliestOccurrence          $date
	OnlyRelevant          bool
	Sid          int
	OnlyRelevant          bool
}

func scheduleApplyErrataCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleApplyErrata",
		Short: "Schedules an action to apply errata updates to multiple systems.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleApplyErrataFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleApplyErrata)
		},
	}

	cmd.Flags().String("Sids", "", "$desc")
	cmd.Flags().String("ErrataIds", "", "$desc")
	cmd.Flags().String("AllowModules", "", "Allow this API call, despite modular content being present")
	cmd.Flags().String("EarliestOccurrence", "", "")
	cmd.Flags().String("OnlyRelevant", "", "If true not all erratas are applied to all systems. Systems get only the erratas relevant for them.")
	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("OnlyRelevant", "", "")

	return cmd
}

func scheduleApplyErrata(globalFlags *types.GlobalFlags, flags *scheduleApplyErrataFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sids, flags.ErrataIds, flags.AllowModules, flags.EarliestOccurrence, flags.OnlyRelevant, flags.Sid, flags.OnlyRelevant)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

