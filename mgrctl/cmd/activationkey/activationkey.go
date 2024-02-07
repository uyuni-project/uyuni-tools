
package activationkey

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activationkey",
		Short: "Contains methods to access common activation key functions
 available from the web interface.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(addPackagesCommand(globalFlags))
	cmd.AddCommand(addEntitlementsCommand(globalFlags))
	cmd.AddCommand(removeServerGroupsCommand(globalFlags))
	cmd.AddCommand(removeConfigChannelsCommand(globalFlags))
	cmd.AddCommand(listActivationKeysCommand(globalFlags))
	cmd.AddCommand(disableConfigDeploymentCommand(globalFlags))
	cmd.AddCommand(removePackagesCommand(globalFlags))
	cmd.AddCommand(checkConfigDeploymentCommand(globalFlags))
	cmd.AddCommand(setConfigChannelsCommand(globalFlags))
	cmd.AddCommand(listActivatedSystemsCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))
	cmd.AddCommand(listConfigChannelsCommand(globalFlags))
	cmd.AddCommand(addConfigChannelsCommand(globalFlags))
	cmd.AddCommand(listChannelsCommand(globalFlags))
	cmd.AddCommand(removeChildChannelsCommand(globalFlags))
	cmd.AddCommand(setDetailsCommand(globalFlags))
	cmd.AddCommand(cloneCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(enableConfigDeploymentCommand(globalFlags))
	cmd.AddCommand(removeEntitlementsCommand(globalFlags))
	cmd.AddCommand(addChildChannelsCommand(globalFlags))
	cmd.AddCommand(addServerGroupsCommand(globalFlags))

	return cmd
}
