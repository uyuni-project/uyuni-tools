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

type scheduleRebootFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	EarliestOccurrence          $date
}

func scheduleRebootCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleReboot",
		Short: "Schedule a reboot for a system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleRebootFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleReboot)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("EarliestOccurrence", "", "")

	return cmd
}

func scheduleReboot(globalFlags *types.GlobalFlags, flags *scheduleRebootFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.EarliestOccurrence)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

