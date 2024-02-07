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

type listAllArchivedActionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listAllArchivedActionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAllArchivedActions",
		Short: "Returns a list of actions that have been archived.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAllArchivedActionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAllArchivedActions)
		},
	}


	return cmd
}

func listAllArchivedActions(globalFlags *types.GlobalFlags, flags *listAllArchivedActionsFlags, cmd *cobra.Command, args []string) error {

res, err := schedule.Schedule(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

