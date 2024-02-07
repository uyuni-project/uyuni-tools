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

type listPhysicalSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listPhysicalSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listPhysicalSystems",
		Short: "Returns a list of all Physical servers visible to the user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listPhysicalSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listPhysicalSystems)
		},
	}

	return cmd
}

func listPhysicalSystems(globalFlags *types.GlobalFlags, flags *listPhysicalSystemsFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
