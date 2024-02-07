package software

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "software",
		Short: "Provides methods to access and modify many aspects of a channel.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(updateRepoCommand(globalFlags))
	cmd.AddCommand(addPackagesCommand(globalFlags))
	cmd.AddCommand(removeRepoFilterCommand(globalFlags))
	cmd.AddCommand(listErrataCommand(globalFlags))
	cmd.AddCommand(updateRepoUrlCommand(globalFlags))
	cmd.AddCommand(removeErrataCommand(globalFlags))
	cmd.AddCommand(setRepoFiltersCommand(globalFlags))
	cmd.AddCommand(associateRepoCommand(globalFlags))
	cmd.AddCommand(createRepoCommand(globalFlags))
	cmd.AddCommand(setDetailsCommand(globalFlags))
	cmd.AddCommand(listSystemChannelsCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(mergePackagesCommand(globalFlags))
	cmd.AddCommand(mergeErrataCommand(globalFlags))
	cmd.AddCommand(listUserReposCommand(globalFlags))
	cmd.AddCommand(addRepoFilterCommand(globalFlags))
	cmd.AddCommand(updateRepoSslCommand(globalFlags))
	cmd.AddCommand(applyChannelStateCommand(globalFlags))
	cmd.AddCommand(listArchesCommand(globalFlags))
	cmd.AddCommand(removePackagesCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(listLatestPackagesCommand(globalFlags))
	cmd.AddCommand(getRepoSyncCronExpressionCommand(globalFlags))
	cmd.AddCommand(alignMetadataCommand(globalFlags))
	cmd.AddCommand(disassociateRepoCommand(globalFlags))
	cmd.AddCommand(listChildrenCommand(globalFlags))
	cmd.AddCommand(regenerateNeededCacheCommand(globalFlags))
	cmd.AddCommand(listErrataByTypeCommand(globalFlags))
	cmd.AddCommand(listPackagesWithoutChannelCommand(globalFlags))
	cmd.AddCommand(listErrataNeedingSyncCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))
	cmd.AddCommand(clearRepoFiltersCommand(globalFlags))
	cmd.AddCommand(listSubscribedSystemsCommand(globalFlags))
	cmd.AddCommand(listRepoFiltersCommand(globalFlags))
	cmd.AddCommand(getChannelLastBuildByIdCommand(globalFlags))
	cmd.AddCommand(isUserSubscribableCommand(globalFlags))
	cmd.AddCommand(setUserManageableCommand(globalFlags))
	cmd.AddCommand(listAllPackagesCommand(globalFlags))
	cmd.AddCommand(regenerateYumCacheCommand(globalFlags))
	cmd.AddCommand(syncErrataCommand(globalFlags))
	cmd.AddCommand(isUserManageableCommand(globalFlags))
	cmd.AddCommand(isExistingCommand(globalFlags))
	cmd.AddCommand(syncRepoCommand(globalFlags))
	cmd.AddCommand(setGloballySubscribableCommand(globalFlags))
	cmd.AddCommand(isGloballySubscribableCommand(globalFlags))
	cmd.AddCommand(removeRepoCommand(globalFlags))
	cmd.AddCommand(updateRepoLabelCommand(globalFlags))
	cmd.AddCommand(getRepoDetailsCommand(globalFlags))
	cmd.AddCommand(cloneCommand(globalFlags))
	cmd.AddCommand(setContactDetailsCommand(globalFlags))
	cmd.AddCommand(listChannelReposCommand(globalFlags))
	cmd.AddCommand(setUserSubscribableCommand(globalFlags))

	return cmd
}
