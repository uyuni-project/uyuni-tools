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

type removeAssignedSystemGroupsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login          string
	$param.getFlagName()          $param.getType()
	SetDefault          bool
}

func removeAssignedSystemGroupsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeAssignedSystemGroups",
		Short: "Remove system groups from a user's list of assigned system groups.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeAssignedSystemGroupsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeAssignedSystemGroups)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("SetDefault", "", "Should system groups also be removed from the user's list of default system groups.")

	return cmd
}

func removeAssignedSystemGroups(globalFlags *types.GlobalFlags, flags *removeAssignedSystemGroupsFlags, cmd *cobra.Command, args []string) error {

res, err := user.User(&flags.ConnectionDetails, flags.Login, flags.$param.getFlagName(), flags.SetDefault)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

