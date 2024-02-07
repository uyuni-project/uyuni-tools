
package actionchain

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "actionchain",
		Short: "Provides the namespace for the Action Chain methods.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(removeActionCommand(globalFlags))
	cmd.AddCommand(createChainCommand(globalFlags))
	cmd.AddCommand(deleteChainCommand(globalFlags))
	cmd.AddCommand(addErrataUpdateCommand(globalFlags))
	cmd.AddCommand(addPackageRemovalCommand(globalFlags))
	cmd.AddCommand(renameChainCommand(globalFlags))
	cmd.AddCommand(addScriptRunCommand(globalFlags))
	cmd.AddCommand(listChainsCommand(globalFlags))
	cmd.AddCommand(addPackageInstallCommand(globalFlags))
	cmd.AddCommand(listChainActionsCommand(globalFlags))
	cmd.AddCommand(addPackageUpgradeCommand(globalFlags))
	cmd.AddCommand(addConfigurationDeploymentCommand(globalFlags))
	cmd.AddCommand(addPackageVerifyCommand(globalFlags))
	cmd.AddCommand(scheduleChainCommand(globalFlags))
	cmd.AddCommand(addSystemRebootCommand(globalFlags))

	return cmd
}
