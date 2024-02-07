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

type updateCalendarFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	$param.getFlagName()          $param.getType()
	$param.getFlagName()          $param.getType()
}

func updateCalendarCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateCalendar",
		Short: "Update a maintenance calendar",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateCalendarFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, updateCalendar)
		},
	}

	cmd.Flags().String("Label", "", "maintenance calendar label")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func updateCalendar(globalFlags *types.GlobalFlags, flags *updateCalendarFlags, cmd *cobra.Command, args []string) error {

res, err := maintenance.Maintenance(&flags.ConnectionDetails, flags.Label, flags.$param.getFlagName(), flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

