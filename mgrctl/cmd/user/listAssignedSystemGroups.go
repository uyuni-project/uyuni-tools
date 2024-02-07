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

type listAssignedSystemGroupsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login                 string
}

func listAssignedSystemGroupsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAssignedSystemGroups",
		Short: "Returns the system groups that a user can administer.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAssignedSystemGroupsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAssignedSystemGroups)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")

	return cmd
}

func listAssignedSystemGroups(globalFlags *types.GlobalFlags, flags *listAssignedSystemGroupsFlags, cmd *cobra.Command, args []string) error {

	res, err := user.User(&flags.ConnectionDetails, flags.Login)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
