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

type createScheduleFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
	Type          string
	Calendar          string
}

func createScheduleCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createSchedule",
		Short: "Create a new maintenance Schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createScheduleFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createSchedule)
		},
	}

	cmd.Flags().String("Name", "", "maintenance schedule name")
	cmd.Flags().String("Type", "", "schedule type: single, multi")
	cmd.Flags().String("Calendar", "", "maintenance calendar label")

	return cmd
}

func createSchedule(globalFlags *types.GlobalFlags, flags *createScheduleFlags, cmd *cobra.Command, args []string) error {

res, err := maintenance.Maintenance(&flags.ConnectionDetails, flags.Name, flags.Type, flags.Calendar)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

