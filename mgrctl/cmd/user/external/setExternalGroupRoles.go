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

type setExternalGroupRolesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
}

func setExternalGroupRolesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setExternalGroupRoles",
		Short: "Update the roles for an external group. Replace previously set roles
 with the ones passed in here. Can only be called by a #product() Administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setExternalGroupRolesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setExternalGroupRoles)
		},
	}

	cmd.Flags().String("Name", "", "Name of the external group.")

	return cmd
}

func setExternalGroupRoles(globalFlags *types.GlobalFlags, flags *setExternalGroupRolesFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

