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

type listFailedActionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listFailedActionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listFailedActions",
		Short: "Returns a list of actions that have failed.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listFailedActionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listFailedActions)
		},
	}

	return cmd
}

func listFailedActions(globalFlags *types.GlobalFlags, flags *listFailedActionsFlags, cmd *cobra.Command, args []string) error {

	res, err := schedule.Schedule(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
