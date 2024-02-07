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

type refreshCalendarFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	$param.getFlagName()          $param.getType()
}

func refreshCalendarCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refreshCalendar",
		Short: "Refresh maintenance calendar data using the configured URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags refreshCalendarFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, refreshCalendar)
		},
	}

	cmd.Flags().String("Label", "", "maintenance calendar label")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func refreshCalendar(globalFlags *types.GlobalFlags, flags *refreshCalendarFlags, cmd *cobra.Command, args []string) error {

res, err := maintenance.Maintenance(&flags.ConnectionDetails, flags.Label, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

