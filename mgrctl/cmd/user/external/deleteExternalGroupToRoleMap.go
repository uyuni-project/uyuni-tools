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

type deleteExternalGroupToRoleMapFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
}

func deleteExternalGroupToRoleMapCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteExternalGroupToRoleMap",
		Short: "Delete the role map for an external group. Can only be called
 by a #product() Administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteExternalGroupToRoleMapFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteExternalGroupToRoleMap)
		},
	}

	cmd.Flags().String("Name", "", "Name of the external group.")

	return cmd
}

func deleteExternalGroupToRoleMap(globalFlags *types.GlobalFlags, flags *deleteExternalGroupToRoleMapFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

