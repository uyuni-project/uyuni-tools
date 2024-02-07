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

type setCreateDefaultSystemGroupFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	CreateDefaultSystemGroup          bool
}

func setCreateDefaultSystemGroupCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setCreateDefaultSystemGroup",
		Short: "Sets the value of the createDefaultSystemGroup setting.
 If True this will cause there to be a system group created (with the same name
 as the user) every time a new user is created, with the user automatically given
 permission to that system group and the system group being set as the default
 group for the user (so every time the user registers a system it will be
 placed in that system group by default). This can be useful if different
 users will administer different groups of servers in the same organization.
 Can only be called by an org_admin.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setCreateDefaultSystemGroupFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setCreateDefaultSystemGroup)
		},
	}

	cmd.Flags().String("CreateDefaultSystemGroup", "", "true if we should automatically create system groups, false otherwise.")

	return cmd
}

func setCreateDefaultSystemGroup(globalFlags *types.GlobalFlags, flags *setCreateDefaultSystemGroupFlags, cmd *cobra.Command, args []string) error {

res, err := user.User(&flags.ConnectionDetails, flags.CreateDefaultSystemGroup)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

