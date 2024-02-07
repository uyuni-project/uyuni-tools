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

type listAssignableRolesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listAssignableRolesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAssignableRoles",
		Short: "Returns a list of user roles that this user can assign to others.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAssignableRolesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAssignableRoles)
		},
	}


	return cmd
}

func listAssignableRoles(globalFlags *types.GlobalFlags, flags *listAssignableRolesFlags, cmd *cobra.Command, args []string) error {

res, err := user.User(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

