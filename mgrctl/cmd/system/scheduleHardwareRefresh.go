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

type scheduleHardwareRefreshFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	EarliestOccurrence          $date
}

func scheduleHardwareRefreshCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleHardwareRefresh",
		Short: "Schedule a hardware refresh for a system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleHardwareRefreshFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleHardwareRefresh)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("EarliestOccurrence", "", "")

	return cmd
}

func scheduleHardwareRefresh(globalFlags *types.GlobalFlags, flags *scheduleHardwareRefreshFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.EarliestOccurrence)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

