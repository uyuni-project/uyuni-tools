
package configchannel

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configchannel",
		Short: "Provides methods to access and modify many aspects of
 configuration channels.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(lookupChannelInfoCommand(globalFlags))
	cmd.AddCommand(listGlobalsCommand(globalFlags))
	cmd.AddCommand(deleteFileRevisionsCommand(globalFlags))
	cmd.AddCommand(channelExistsCommand(globalFlags))
	cmd.AddCommand(updateCommand(globalFlags))
	cmd.AddCommand(scheduleFileComparisonsCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(createOrUpdatePathCommand(globalFlags))
	cmd.AddCommand(deployAllSystemsCommand(globalFlags))
	cmd.AddCommand(listSubscribedSystemsCommand(globalFlags))
	cmd.AddCommand(deleteFilesCommand(globalFlags))
	cmd.AddCommand(updateInitSlsCommand(globalFlags))
	cmd.AddCommand(getEncodedFileRevisionCommand(globalFlags))
	cmd.AddCommand(lookupFileInfoCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(syncSaltFilesOnDiskCommand(globalFlags))
	cmd.AddCommand(getFileRevisionsCommand(globalFlags))
	cmd.AddCommand(listFilesCommand(globalFlags))
	cmd.AddCommand(deleteChannelsCommand(globalFlags))
	cmd.AddCommand(createOrUpdateSymlinkCommand(globalFlags))
	cmd.AddCommand(listAssignedSystemGroupsCommand(globalFlags))
	cmd.AddCommand(getFileRevisionCommand(globalFlags))

	return cmd
}
