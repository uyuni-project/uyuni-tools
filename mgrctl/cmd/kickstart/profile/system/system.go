
package system

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "Provides methods to set various properties of a kickstart profile.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(setPartitioningSchemeCommand(globalFlags))
	cmd.AddCommand(getPartitioningSchemeCommand(globalFlags))
	cmd.AddCommand(addFilePreservationsCommand(globalFlags))
	cmd.AddCommand(listKeysCommand(globalFlags))
	cmd.AddCommand(disableRemoteCommandsCommand(globalFlags))
	cmd.AddCommand(disableConfigManagementCommand(globalFlags))
	cmd.AddCommand(setSELinuxCommand(globalFlags))
	cmd.AddCommand(enableRemoteCommandsCommand(globalFlags))
	cmd.AddCommand(setRegistrationTypeCommand(globalFlags))
	cmd.AddCommand(removeKeysCommand(globalFlags))
	cmd.AddCommand(removeFilePreservationsCommand(globalFlags))
	cmd.AddCommand(checkConfigManagementCommand(globalFlags))
	cmd.AddCommand(addKeysCommand(globalFlags))
	cmd.AddCommand(getSELinuxCommand(globalFlags))
	cmd.AddCommand(getRegistrationTypeCommand(globalFlags))
	cmd.AddCommand(listFilePreservationsCommand(globalFlags))
	cmd.AddCommand(enableConfigManagementCommand(globalFlags))
	cmd.AddCommand(getLocaleCommand(globalFlags))
	cmd.AddCommand(checkRemoteCommandsCommand(globalFlags))
	cmd.AddCommand(setLocaleCommand(globalFlags))

	return cmd
}
