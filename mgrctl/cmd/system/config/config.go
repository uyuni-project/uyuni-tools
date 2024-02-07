
package config

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Provides methods to access and modify many aspects of
 configuration channels and server association.
 basically system.config name space",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(addChannelsCommand(globalFlags))
	cmd.AddCommand(deployAllCommand(globalFlags))
	cmd.AddCommand(deleteFilesCommand(globalFlags))
	cmd.AddCommand(listChannelsCommand(globalFlags))
	cmd.AddCommand(setChannelsCommand(globalFlags))
	cmd.AddCommand(scheduleApplyConfigChannelCommand(globalFlags))
	cmd.AddCommand(lookupFileInfoCommand(globalFlags))
	cmd.AddCommand(createOrUpdatePathCommand(globalFlags))
	cmd.AddCommand(listFilesCommand(globalFlags))
	cmd.AddCommand(createOrUpdateSymlinkCommand(globalFlags))
	cmd.AddCommand(removeChannelsCommand(globalFlags))

	return cmd
}
