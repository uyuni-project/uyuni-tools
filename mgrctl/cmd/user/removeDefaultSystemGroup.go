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

type removeDefaultSystemGroupFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login          string
	SgName          string
}

func removeDefaultSystemGroupCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeDefaultSystemGroup",
		Short: "Remove a system group from user's list of default system groups.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeDefaultSystemGroupFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeDefaultSystemGroup)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")
	cmd.Flags().String("SgName", "", "server group name")

	return cmd
}

func removeDefaultSystemGroup(globalFlags *types.GlobalFlags, flags *removeDefaultSystemGroupFlags, cmd *cobra.Command, args []string) error {

res, err := user.User(&flags.ConnectionDetails, flags.Login, flags.SgName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

