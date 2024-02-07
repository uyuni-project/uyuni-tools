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

type addRoleFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login          string
	Role          string
}

func addRoleCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addRole",
		Short: "Adds a role to a user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addRoleFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addRole)
		},
	}

	cmd.Flags().String("Login", "", "User login name to update.")
	cmd.Flags().String("Role", "", "Role label to add.  Can be any of: satellite_admin, org_admin, channel_admin, config_admin, system_group_admin, or activation_key_admin.")

	return cmd
}

func addRole(globalFlags *types.GlobalFlags, flags *addRoleFlags, cmd *cobra.Command, args []string) error {

res, err := user.User(&flags.ConnectionDetails, flags.Login, flags.Role)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

