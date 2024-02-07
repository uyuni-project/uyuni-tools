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

type listInProgressSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ActionId              int
}

func listInProgressSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listInProgressSystems",
		Short: "Returns a list of systems that have a specific action in progress.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listInProgressSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listInProgressSystems)
		},
	}

	cmd.Flags().String("ActionId", "", "")

	return cmd
}

func listInProgressSystems(globalFlags *types.GlobalFlags, flags *listInProgressSystemsFlags, cmd *cobra.Command, args []string) error {

	res, err := schedule.Schedule(&flags.ConnectionDetails, flags.ActionId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
