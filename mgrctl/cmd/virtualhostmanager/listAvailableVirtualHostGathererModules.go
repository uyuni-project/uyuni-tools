package virtualhostmanager

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/virtualhostmanager"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listAvailableVirtualHostGathererModulesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listAvailableVirtualHostGathererModulesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAvailableVirtualHostGathererModules",
		Short: "List all available modules from virtual-host-gatherer",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAvailableVirtualHostGathererModulesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAvailableVirtualHostGathererModules)
		},
	}


	return cmd
}

func listAvailableVirtualHostGathererModules(globalFlags *types.GlobalFlags, flags *listAvailableVirtualHostGathererModulesFlags, cmd *cobra.Command, args []string) error {

res, err := virtualhostmanager.Virtualhostmanager(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

