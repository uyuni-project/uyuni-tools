
package systemgroup

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "systemgroup",
		Short: "Provides methods to access and modify system groups.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listSystemsMinimalCommand(globalFlags))
	cmd.AddCommand(addOrRemoveAdminsCommand(globalFlags))
	cmd.AddCommand(updateCommand(globalFlags))
	cmd.AddCommand(listActiveSystemsInGroupCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(listInactiveSystemsInGroupCommand(globalFlags))
	cmd.AddCommand(listSystemsCommand(globalFlags))
	cmd.AddCommand(listAssignedConfigChannelsCommand(globalFlags))
	cmd.AddCommand(listAdministratorsCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))
	cmd.AddCommand(listGroupsWithNoAssociatedAdminsCommand(globalFlags))
	cmd.AddCommand(subscribeConfigChannelCommand(globalFlags))
	cmd.AddCommand(addOrRemoveSystemsCommand(globalFlags))
	cmd.AddCommand(listAssignedFormualsCommand(globalFlags))
	cmd.AddCommand(unsubscribeConfigChannelCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(scheduleApplyErrataToActiveCommand(globalFlags))
	cmd.AddCommand(listAllGroupsCommand(globalFlags))

	return cmd
}
