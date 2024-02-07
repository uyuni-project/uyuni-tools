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

type assignScheduleToSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ScheduleName          string
	$param.getFlagName()          $param.getType()
	$param.getFlagName()          $param.getType()
}

func assignScheduleToSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assignScheduleToSystems",
		Short: "Assign schedule with given name to systems with given IDs.
 Throws a PermissionCheckFailureException when some of the systems are not accessible by the user.
 Throws a InvalidParameterException when some of the systems have pending actions that are not allowed in the
 maintenance mode.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags assignScheduleToSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, assignScheduleToSystems)
		},
	}

	cmd.Flags().String("ScheduleName", "", "The schedule name")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func assignScheduleToSystems(globalFlags *types.GlobalFlags, flags *assignScheduleToSystemsFlags, cmd *cobra.Command, args []string) error {

res, err := maintenance.Maintenance(&flags.ConnectionDetails, flags.ScheduleName, flags.$param.getFlagName(), flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

