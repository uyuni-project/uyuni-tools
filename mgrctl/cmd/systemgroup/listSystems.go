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

type listSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemGroupName          string
}

func listSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSystems",
		Short: "Return a list of systems associated with this system group.
 User must have access to this system group.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSystems)
		},
	}

	cmd.Flags().String("SystemGroupName", "", "")

	return cmd
}

func listSystems(globalFlags *types.GlobalFlags, flags *listSystemsFlags, cmd *cobra.Command, args []string) error {

res, err := systemgroup.Systemgroup(&flags.ConnectionDetails, flags.SystemGroupName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

