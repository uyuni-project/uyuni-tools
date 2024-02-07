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

type listAllActionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listAllActionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAllActions",
		Short: "Returns a list of all actions.  This includes completed, in progress,
 failed and archived actions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAllActionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAllActions)
		},
	}


	return cmd
}

func listAllActions(globalFlags *types.GlobalFlags, flags *listAllActionsFlags, cmd *cobra.Command, args []string) error {

res, err := schedule.Schedule(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

