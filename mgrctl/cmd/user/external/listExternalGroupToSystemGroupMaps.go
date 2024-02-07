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

type listExternalGroupToSystemGroupMapsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listExternalGroupToSystemGroupMapsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listExternalGroupToSystemGroupMaps",
		Short: "List server group mappings for all known external groups. Can only be
 called by an org_admin.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listExternalGroupToSystemGroupMapsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listExternalGroupToSystemGroupMaps)
		},
	}


	return cmd
}

func listExternalGroupToSystemGroupMaps(globalFlags *types.GlobalFlags, flags *listExternalGroupToSystemGroupMapsFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

