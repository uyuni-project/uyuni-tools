package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listInactiveSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Days          int
}

func listInactiveSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listInactiveSystems",
		Short: "Lists systems that have been inactive for the default period of
          inactivity",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listInactiveSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listInactiveSystems)
		},
	}

	cmd.Flags().String("Days", "", "")

	return cmd
}

func listInactiveSystems(globalFlags *types.GlobalFlags, flags *listInactiveSystemsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Days)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

