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

type getExternalGroupToRoleMapFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
}

func getExternalGroupToRoleMapCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getExternalGroupToRoleMap",
		Short: "Get a representation of the role mapping for an external group.
 Can only be called by a #product() Administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getExternalGroupToRoleMapFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getExternalGroupToRoleMap)
		},
	}

	cmd.Flags().String("Name", "", "Name of the external group.")

	return cmd
}

func getExternalGroupToRoleMap(globalFlags *types.GlobalFlags, flags *getExternalGroupToRoleMapFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

