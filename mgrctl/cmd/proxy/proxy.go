
package proxy

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "Provides methods to activate/deactivate a proxy
 server.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listAvailableProxyChannelsCommand(globalFlags))
	cmd.AddCommand(containerConfigCommand(globalFlags))
	cmd.AddCommand(createMonitoringScoutCommand(globalFlags))
	cmd.AddCommand(activateProxyCommand(globalFlags))
	cmd.AddCommand(listProxyClientsCommand(globalFlags))
	cmd.AddCommand(isProxyCommand(globalFlags))
	cmd.AddCommand(listProxiesCommand(globalFlags))
	cmd.AddCommand(deactivateProxyCommand(globalFlags))

	return cmd
}
