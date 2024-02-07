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

type scheduleApplyHighstateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	EarliestOccurrence          $date
	Test          bool
	Sids          []int
}

func scheduleApplyHighstateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleApplyHighstate",
		Short: "Schedule highstate application for a given system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleApplyHighstateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleApplyHighstate)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("EarliestOccurrence", "", "")
	cmd.Flags().String("Test", "", "Run states in test-only mode")
	cmd.Flags().String("Sids", "", "$desc")

	return cmd
}

func scheduleApplyHighstate(globalFlags *types.GlobalFlags, flags *scheduleApplyHighstateFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.EarliestOccurrence, flags.Test, flags.Sids)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

