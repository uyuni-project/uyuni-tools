
package profile

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Provides methods to access and modify many aspects of
 a kickstart profile.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(getCfgPreservationCommand(globalFlags))
	cmd.AddCommand(getVariablesCommand(globalFlags))
	cmd.AddCommand(getUpdateTypeCommand(globalFlags))
	cmd.AddCommand(setChildChannelsCommand(globalFlags))
	cmd.AddCommand(downloadRenderedKickstartCommand(globalFlags))
	cmd.AddCommand(setAdvancedOptionsCommand(globalFlags))
	cmd.AddCommand(setCustomOptionsCommand(globalFlags))
	cmd.AddCommand(setVirtualizationTypeCommand(globalFlags))
	cmd.AddCommand(comparePackagesCommand(globalFlags))
	cmd.AddCommand(getAvailableRepositoriesCommand(globalFlags))
	cmd.AddCommand(addIpRangeCommand(globalFlags))
	cmd.AddCommand(addScriptCommand(globalFlags))
	cmd.AddCommand(setRepositoriesCommand(globalFlags))
	cmd.AddCommand(compareAdvancedOptionsCommand(globalFlags))
	cmd.AddCommand(setKickstartTreeCommand(globalFlags))
	cmd.AddCommand(compareActivationKeysCommand(globalFlags))
	cmd.AddCommand(getAdvancedOptionsCommand(globalFlags))
	cmd.AddCommand(removeScriptCommand(globalFlags))
	cmd.AddCommand(getKickstartTreeCommand(globalFlags))
	cmd.AddCommand(setCfgPreservationCommand(globalFlags))
	cmd.AddCommand(setVariablesCommand(globalFlags))
	cmd.AddCommand(setUpdateTypeCommand(globalFlags))
	cmd.AddCommand(getVirtualizationTypeCommand(globalFlags))
	cmd.AddCommand(listScriptsCommand(globalFlags))
	cmd.AddCommand(downloadKickstartCommand(globalFlags))
	cmd.AddCommand(getChildChannelsCommand(globalFlags))
	cmd.AddCommand(getRepositoriesCommand(globalFlags))
	cmd.AddCommand(removeIpRangeCommand(globalFlags))
	cmd.AddCommand(setLoggingCommand(globalFlags))
	cmd.AddCommand(orderScriptsCommand(globalFlags))
	cmd.AddCommand(listIpRangesCommand(globalFlags))
	cmd.AddCommand(getCustomOptionsCommand(globalFlags))

	return cmd
}
