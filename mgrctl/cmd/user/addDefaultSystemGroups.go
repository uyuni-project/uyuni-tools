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

type addDefaultSystemGroupsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login          string
	$param.getFlagName()          $param.getType()
}

func addDefaultSystemGroupsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addDefaultSystemGroups",
		Short: "Add system groups to user's list of default system groups.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addDefaultSystemGroupsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addDefaultSystemGroups)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func addDefaultSystemGroups(globalFlags *types.GlobalFlags, flags *addDefaultSystemGroupsFlags, cmd *cobra.Command, args []string) error {

res, err := user.User(&flags.ConnectionDetails, flags.Login, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

