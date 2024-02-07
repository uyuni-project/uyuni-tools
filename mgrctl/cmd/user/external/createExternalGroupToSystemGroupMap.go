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

type createExternalGroupToSystemGroupMapFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
}

func createExternalGroupToSystemGroupMapCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createExternalGroupToSystemGroupMap",
		Short: "Externally authenticated users may be members of external groups. You
 can use these groups to give access to server groups to the users when they log in.
 Can only be called by an org_admin.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createExternalGroupToSystemGroupMapFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createExternalGroupToSystemGroupMap)
		},
	}

	cmd.Flags().String("Name", "", "Name of the external group. Must be unique.")

	return cmd
}

func createExternalGroupToSystemGroupMap(globalFlags *types.GlobalFlags, flags *createExternalGroupToSystemGroupMapFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

