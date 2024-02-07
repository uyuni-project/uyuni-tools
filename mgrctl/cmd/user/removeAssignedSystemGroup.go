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

type removeAssignedSystemGroupFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login          string
	SgName          string
	SetDefault          bool
}

func removeAssignedSystemGroupCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeAssignedSystemGroup",
		Short: "Remove system group from the user's list of assigned system groups.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeAssignedSystemGroupFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeAssignedSystemGroup)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")
	cmd.Flags().String("SgName", "", "server group name")
	cmd.Flags().String("SetDefault", "", "Should system group also be removed from the user's list of default system groups.")

	return cmd
}

func removeAssignedSystemGroup(globalFlags *types.GlobalFlags, flags *removeAssignedSystemGroupFlags, cmd *cobra.Command, args []string) error {

res, err := user.User(&flags.ConnectionDetails, flags.Login, flags.SgName, flags.SetDefault)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

