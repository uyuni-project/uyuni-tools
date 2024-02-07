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

type addAssignedSystemGroupFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login                 string
	SgName                string
	SetDefault            bool
}

func addAssignedSystemGroupCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addAssignedSystemGroup",
		Short: "Add system group to user's list of assigned system groups.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addAssignedSystemGroupFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addAssignedSystemGroup)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")
	cmd.Flags().String("SgName", "", "")
	cmd.Flags().String("SetDefault", "", "Should system group also be added to user's list of default system groups.")

	return cmd
}

func addAssignedSystemGroup(globalFlags *types.GlobalFlags, flags *addAssignedSystemGroupFlags, cmd *cobra.Command, args []string) error {

	res, err := user.User(&flags.ConnectionDetails, flags.Login, flags.SgName, flags.SetDefault)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
