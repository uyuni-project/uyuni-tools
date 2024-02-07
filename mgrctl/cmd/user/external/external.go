
package external

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "external",
		Short: "If you are using IPA integration to allow authentication of users from
 an external IPA server (rare) the users will still need to be created in the #product()
 database. Methods in this namespace allow you to configure some specifics of how this
 happens, like what organization they are created in or what roles they will have.
 These options can also be set in the web admin interface.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(getExternalGroupToSystemGroupMapCommand(globalFlags))
	cmd.AddCommand(createExternalGroupToSystemGroupMapCommand(globalFlags))
	cmd.AddCommand(deleteExternalGroupToRoleMapCommand(globalFlags))
	cmd.AddCommand(setKeepTemporaryRolesCommand(globalFlags))
	cmd.AddCommand(listExternalGroupToSystemGroupMapsCommand(globalFlags))
	cmd.AddCommand(getDefaultOrgCommand(globalFlags))
	cmd.AddCommand(getExternalGroupToRoleMapCommand(globalFlags))
	cmd.AddCommand(getKeepTemporaryRolesCommand(globalFlags))
	cmd.AddCommand(setExternalGroupRolesCommand(globalFlags))
	cmd.AddCommand(listExternalGroupToRoleMapsCommand(globalFlags))
	cmd.AddCommand(deleteExternalGroupToSystemGroupMapCommand(globalFlags))
	cmd.AddCommand(setUseOrgUnitCommand(globalFlags))
	cmd.AddCommand(createExternalGroupToRoleMapCommand(globalFlags))
	cmd.AddCommand(getUseOrgUnitCommand(globalFlags))
	cmd.AddCommand(setExternalGroupSystemGroupsCommand(globalFlags))
	cmd.AddCommand(setDefaultOrgCommand(globalFlags))

	return cmd
}
