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

type addAssignedSystemGroupsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login          string
	$param.getFlagName()          $param.getType()
	SetDefault          bool
}

func addAssignedSystemGroupsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addAssignedSystemGroups",
		Short: "Add system groups to user's list of assigned system groups.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addAssignedSystemGroupsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addAssignedSystemGroups)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("SetDefault", "", "Should system groups also be added to user's list of default system groups.")

	return cmd
}

func addAssignedSystemGroups(globalFlags *types.GlobalFlags, flags *addAssignedSystemGroupsFlags, cmd *cobra.Command, args []string) error {

res, err := user.User(&flags.ConnectionDetails, flags.Login, flags.$param.getFlagName(), flags.SetDefault)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

