package user

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/user"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type removeRoleFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login          string
	Role          string
}

func removeRoleCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeRole",
		Short: "Remove a role from a user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeRoleFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeRole)
		},
	}

	cmd.Flags().String("Login", "", "User login name to update.")
	cmd.Flags().String("Role", "", "Role label to remove.  Can be any of: satellite_admin, org_admin, channel_admin, config_admin, system_group_admin, or activation_key_admin.")

	return cmd
}

func removeRole(globalFlags *types.GlobalFlags, flags *removeRoleFlags, cmd *cobra.Command, args []string) error {

res, err := user.User(&flags.ConnectionDetails, flags.Login, flags.Role)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

