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

type listCompletedSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ActionId              int
}

func listCompletedSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listCompletedSystems",
		Short: "Returns a list of systems that have completed a specific action.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listCompletedSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listCompletedSystems)
		},
	}

	cmd.Flags().String("ActionId", "", "")

	return cmd
}

func listCompletedSystems(globalFlags *types.GlobalFlags, flags *listCompletedSystemsFlags, cmd *cobra.Command, args []string) error {

	res, err := schedule.Schedule(&flags.ConnectionDetails, flags.ActionId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
