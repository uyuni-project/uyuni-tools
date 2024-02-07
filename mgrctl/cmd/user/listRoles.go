package user

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/api/user"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listRolesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login                 string
}

func listRolesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listRoles",
		Short: "Returns a list of the user's roles.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listRolesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listRoles)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")

	return cmd
}

func listRoles(globalFlags *types.GlobalFlags, flags *listRolesFlags, cmd *cobra.Command, args []string) error {

	res, err := user.User(&flags.ConnectionDetails, flags.Login)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
