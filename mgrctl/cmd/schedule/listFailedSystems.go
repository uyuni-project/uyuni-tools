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

type listFailedSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ActionId          int
}

func listFailedSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listFailedSystems",
		Short: "Returns a list of systems that have failed a specific action.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listFailedSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listFailedSystems)
		},
	}

	cmd.Flags().String("ActionId", "", "")

	return cmd
}

func listFailedSystems(globalFlags *types.GlobalFlags, flags *listFailedSystemsFlags, cmd *cobra.Command, args []string) error {

res, err := schedule.Schedule(&flags.ConnectionDetails, flags.ActionId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

