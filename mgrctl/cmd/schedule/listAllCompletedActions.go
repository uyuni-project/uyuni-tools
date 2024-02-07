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

type listAllCompletedActionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listAllCompletedActionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAllCompletedActions",
		Short: "Returns a list of actions that have been completed.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAllCompletedActionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAllCompletedActions)
		},
	}


	return cmd
}

func listAllCompletedActions(globalFlags *types.GlobalFlags, flags *listAllCompletedActionsFlags, cmd *cobra.Command, args []string) error {

res, err := schedule.Schedule(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

