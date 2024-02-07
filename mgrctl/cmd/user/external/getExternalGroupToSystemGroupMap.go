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

type getExternalGroupToSystemGroupMapFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
}

func getExternalGroupToSystemGroupMapCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getExternalGroupToSystemGroupMap",
		Short: "Get a representation of the server group mapping for an external
 group. Can only be called by an org_admin.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getExternalGroupToSystemGroupMapFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getExternalGroupToSystemGroupMap)
		},
	}

	cmd.Flags().String("Name", "", "Name of the external group.")

	return cmd
}

func getExternalGroupToSystemGroupMap(globalFlags *types.GlobalFlags, flags *getExternalGroupToSystemGroupMapFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

