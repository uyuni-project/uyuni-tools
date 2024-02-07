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

type listArchivedActionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listArchivedActionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listArchivedActions",
		Short: "Returns a list of actions that have been archived.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listArchivedActionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listArchivedActions)
		},
	}

	return cmd
}

func listArchivedActions(globalFlags *types.GlobalFlags, flags *listArchivedActionsFlags, cmd *cobra.Command, args []string) error {

	res, err := schedule.Schedule(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
