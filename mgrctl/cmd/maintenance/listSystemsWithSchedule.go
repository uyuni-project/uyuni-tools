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

type listSystemsWithScheduleFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ScheduleName          string
}

func listSystemsWithScheduleCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSystemsWithSchedule",
		Short: "List IDs of systems that have given schedule assigned
 Throws a PermissionCheckFailureException when some of the systems are not accessible by the user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSystemsWithScheduleFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSystemsWithSchedule)
		},
	}

	cmd.Flags().String("ScheduleName", "", "the schedule name")

	return cmd
}

func listSystemsWithSchedule(globalFlags *types.GlobalFlags, flags *listSystemsWithScheduleFlags, cmd *cobra.Command, args []string) error {

res, err := maintenance.Maintenance(&flags.ConnectionDetails, flags.ScheduleName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

