package systemgroup

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/systemgroup"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listInactiveSystemsInGroupFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemGroupName          string
	DaysInactive          int
}

func listInactiveSystemsInGroupCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listInactiveSystemsInGroup",
		Short: "Lists inactive systems within a server group using a
          specified inactivity time.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listInactiveSystemsInGroupFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listInactiveSystemsInGroup)
		},
	}

	cmd.Flags().String("SystemGroupName", "", "")
	cmd.Flags().String("DaysInactive", "", "Number of days a system           must not check in to be considered inactive.")

	return cmd
}

func listInactiveSystemsInGroup(globalFlags *types.GlobalFlags, flags *listInactiveSystemsInGroupFlags, cmd *cobra.Command, args []string) error {

res, err := systemgroup.Systemgroup(&flags.ConnectionDetails, flags.SystemGroupName, flags.DaysInactive)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

