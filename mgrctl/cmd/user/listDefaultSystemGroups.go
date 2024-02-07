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

type listDefaultSystemGroupsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login                 string
}

func listDefaultSystemGroupsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listDefaultSystemGroups",
		Short: "Returns a user's list of default system groups.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listDefaultSystemGroupsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listDefaultSystemGroups)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")

	return cmd
}

func listDefaultSystemGroups(globalFlags *types.GlobalFlags, flags *listDefaultSystemGroupsFlags, cmd *cobra.Command, args []string) error {

	res, err := user.User(&flags.ConnectionDetails, flags.Login)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
