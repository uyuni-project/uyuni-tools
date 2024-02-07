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

type getUseOrgUnitFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func getUseOrgUnitCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getUseOrgUnit",
		Short: "Get whether we place users into the organization that corresponds
 to the "orgunit" set on the IPA server. The orgunit name must match exactly the
 #product() organization name. Can only be called by a #product() Administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getUseOrgUnitFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getUseOrgUnit)
		},
	}


	return cmd
}

func getUseOrgUnit(globalFlags *types.GlobalFlags, flags *getUseOrgUnitFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

