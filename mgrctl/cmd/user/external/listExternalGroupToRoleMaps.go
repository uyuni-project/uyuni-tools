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

type listExternalGroupToRoleMapsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listExternalGroupToRoleMapsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listExternalGroupToRoleMaps",
		Short: "List role mappings for all known external groups. Can only be called
 by a #product() Administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listExternalGroupToRoleMapsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listExternalGroupToRoleMaps)
		},
	}


	return cmd
}

func listExternalGroupToRoleMaps(globalFlags *types.GlobalFlags, flags *listExternalGroupToRoleMapsFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

