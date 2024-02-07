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

type cancelActionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ActionIds          []int
}

func cancelActionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancelActions",
		Short: "Cancel all actions in given list. If an invalid action is provided,
 none of the actions given will be canceled.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags cancelActionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, cancelActions)
		},
	}

	cmd.Flags().String("ActionIds", "", "$desc")

	return cmd
}

func cancelActions(globalFlags *types.GlobalFlags, flags *cancelActionsFlags, cmd *cobra.Command, args []string) error {

res, err := schedule.Schedule(&flags.ConnectionDetails, flags.ActionIds)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

