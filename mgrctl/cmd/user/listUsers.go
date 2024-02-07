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

type listUsersFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listUsersCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listUsers",
		Short: "Returns a list of users in your organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listUsersFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listUsers)
		},
	}

	return cmd
}

func listUsers(globalFlags *types.GlobalFlags, flags *listUsersFlags, cmd *cobra.Command, args []string) error {

	res, err := user.User(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
