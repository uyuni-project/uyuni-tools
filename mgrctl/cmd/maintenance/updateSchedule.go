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

type updateScheduleFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
	$param.getFlagName()          $param.getType()
	$param.getFlagName()          $param.getType()
}

func updateScheduleCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateSchedule",
		Short: "Update a maintenance schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateScheduleFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, updateSchedule)
		},
	}

	cmd.Flags().String("Name", "", "maintenance schedule name")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func updateSchedule(globalFlags *types.GlobalFlags, flags *updateScheduleFlags, cmd *cobra.Command, args []string) error {

res, err := maintenance.Maintenance(&flags.ConnectionDetails, flags.Name, flags.$param.getFlagName(), flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

