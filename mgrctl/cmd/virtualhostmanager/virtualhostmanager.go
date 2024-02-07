
package virtualhostmanager

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "virtualhostmanager",
		Short: "Provides the namespace for the Virtual Host Manager methods.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(getDetailCommand(globalFlags))
	cmd.AddCommand(listVirtualHostManagersCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(getModuleParametersCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))
	cmd.AddCommand(listAvailableVirtualHostGathererModulesCommand(globalFlags))

	return cmd
}
