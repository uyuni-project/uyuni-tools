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

type listInProgressActionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listInProgressActionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listInProgressActions",
		Short: "Returns a list of actions that are in progress.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listInProgressActionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listInProgressActions)
		},
	}

	return cmd
}

func listInProgressActions(globalFlags *types.GlobalFlags, flags *listInProgressActionsFlags, cmd *cobra.Command, args []string) error {

	res, err := schedule.Schedule(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
