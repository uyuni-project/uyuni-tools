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

type deleteScheduleFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
}

func deleteScheduleCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteSchedule",
		Short: "Remove a maintenance schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteScheduleFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteSchedule)
		},
	}

	cmd.Flags().String("Name", "", "maintenance schedule name")

	return cmd
}

func deleteSchedule(globalFlags *types.GlobalFlags, flags *deleteScheduleFlags, cmd *cobra.Command, args []string) error {

res, err := maintenance.Maintenance(&flags.ConnectionDetails, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

