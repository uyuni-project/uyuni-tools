
package kickstart

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kickstart",
		Short: "Provides methods to create kickstart files",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listAutoinstallableChannelsCommand(globalFlags))
	cmd.AddCommand(createProfileWithCustomUrlCommand(globalFlags))
	cmd.AddCommand(renameProfileCommand(globalFlags))
	cmd.AddCommand(importRawFileCommand(globalFlags))
	cmd.AddCommand(findKickstartForIpCommand(globalFlags))
	cmd.AddCommand(importFileCommand(globalFlags))
	cmd.AddCommand(deleteProfileCommand(globalFlags))
	cmd.AddCommand(cloneProfileCommand(globalFlags))
	cmd.AddCommand(disableProfileCommand(globalFlags))
	cmd.AddCommand(isProfileDisabledCommand(globalFlags))
	cmd.AddCommand(listAllIpRangesCommand(globalFlags))
	cmd.AddCommand(createProfileCommand(globalFlags))
	cmd.AddCommand(listKickstartsCommand(globalFlags))
	cmd.AddCommand(listKickstartableChannelsCommand(globalFlags))

	return cmd
}
