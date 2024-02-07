package schedule

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/schedule"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listCompletedActionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listCompletedActionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listCompletedActions",
		Short: "Returns a list of actions that have completed successfully.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listCompletedActionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listCompletedActions)
		},
	}


	return cmd
}

func listCompletedActions(globalFlags *types.GlobalFlags, flags *listCompletedActionsFlags, cmd *cobra.Command, args []string) error {

res, err := schedule.Schedule(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

