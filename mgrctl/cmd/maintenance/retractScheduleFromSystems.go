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

type retractScheduleFromSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	$param.getFlagName()          $param.getType()
}

func retractScheduleFromSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retractScheduleFromSystems",
		Short: "Retract schedule with given name from systems with given IDs
 Throws a PermissionCheckFailureException when some of the systems are not accessible by the user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags retractScheduleFromSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, retractScheduleFromSystems)
		},
	}

	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func retractScheduleFromSystems(globalFlags *types.GlobalFlags, flags *retractScheduleFromSystemsFlags, cmd *cobra.Command, args []string) error {

res, err := maintenance.Maintenance(&flags.ConnectionDetails, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

