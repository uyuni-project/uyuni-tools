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

type addDefaultSystemGroupFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login                 string
	Name                  string
}

func addDefaultSystemGroupCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addDefaultSystemGroup",
		Short: "Add system group to user's list of default system groups.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addDefaultSystemGroupFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addDefaultSystemGroup)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")
	cmd.Flags().String("Name", "", "server group name")

	return cmd
}

func addDefaultSystemGroup(globalFlags *types.GlobalFlags, flags *addDefaultSystemGroupFlags, cmd *cobra.Command, args []string) error {

	res, err := user.User(&flags.ConnectionDetails, flags.Login, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
