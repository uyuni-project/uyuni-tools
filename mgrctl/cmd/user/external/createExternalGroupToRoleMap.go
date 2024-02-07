package external

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/user/external"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type createExternalGroupToRoleMapFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
}

func createExternalGroupToRoleMapCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createExternalGroupToRoleMap",
		Short: "Externally authenticated users may be members of external groups. You
 can use these groups to assign additional roles to the users when they log in.
 Can only be called by a #product() Administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createExternalGroupToRoleMapFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createExternalGroupToRoleMap)
		},
	}

	cmd.Flags().String("Name", "", "Name of the external group. Must be unique.")

	return cmd
}

func createExternalGroupToRoleMap(globalFlags *types.GlobalFlags, flags *createExternalGroupToRoleMapFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

