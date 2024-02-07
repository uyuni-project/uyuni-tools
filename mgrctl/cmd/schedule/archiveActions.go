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

type archiveActionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ActionIds          []int
}

func archiveActionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archiveActions",
		Short: "Archive all actions in the given list.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags archiveActionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, archiveActions)
		},
	}

	cmd.Flags().String("ActionIds", "", "$desc")

	return cmd
}

func archiveActions(globalFlags *types.GlobalFlags, flags *archiveActionsFlags, cmd *cobra.Command, args []string) error {

res, err := schedule.Schedule(&flags.ConnectionDetails, flags.ActionIds)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

