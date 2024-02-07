
package user

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "User namespace contains methods to access common user functions
 available from the web user interface.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(deleteCommand(globalFlags))
	cmd.AddCommand(addAssignedSystemGroupCommand(globalFlags))
	cmd.AddCommand(enableCommand(globalFlags))
	cmd.AddCommand(setDetailsCommand(globalFlags))
	cmd.AddCommand(removeAssignedSystemGroupsCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(listAssignedSystemGroupsCommand(globalFlags))
	cmd.AddCommand(removeRoleCommand(globalFlags))
	cmd.AddCommand(getCreateDefaultSystemGroupCommand(globalFlags))
	cmd.AddCommand(removeDefaultSystemGroupsCommand(globalFlags))
	cmd.AddCommand(listAssignableRolesCommand(globalFlags))
	cmd.AddCommand(setCreateDefaultSystemGroupCommand(globalFlags))
	cmd.AddCommand(setErrataNotificationsCommand(globalFlags))
	cmd.AddCommand(addAssignedSystemGroupsCommand(globalFlags))
	cmd.AddCommand(removeDefaultSystemGroupCommand(globalFlags))
	cmd.AddCommand(usePamAuthenticationCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(addDefaultSystemGroupCommand(globalFlags))
	cmd.AddCommand(listRolesCommand(globalFlags))
	cmd.AddCommand(addDefaultSystemGroupsCommand(globalFlags))
	cmd.AddCommand(listUsersCommand(globalFlags))
	cmd.AddCommand(disableCommand(globalFlags))
	cmd.AddCommand(listDefaultSystemGroupsCommand(globalFlags))
	cmd.AddCommand(removeAssignedSystemGroupCommand(globalFlags))
	cmd.AddCommand(setReadOnlyCommand(globalFlags))
	cmd.AddCommand(addRoleCommand(globalFlags))

	return cmd
}
