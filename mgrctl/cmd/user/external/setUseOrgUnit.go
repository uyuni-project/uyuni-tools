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

type setUseOrgUnitFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	UseOrgUnit          bool
}

func setUseOrgUnitCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setUseOrgUnit",
		Short: "Set whether we place users into the organization that corresponds
 to the "orgunit" set on the IPA server. The orgunit name must match exactly the
 #product() organization name. Can only be called by a #product() Administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setUseOrgUnitFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setUseOrgUnit)
		},
	}

	cmd.Flags().String("UseOrgUnit", "", "true if we should use the IPA orgunit to determine which organization to create the user in, false otherwise.")

	return cmd
}

func setUseOrgUnit(globalFlags *types.GlobalFlags, flags *setUseOrgUnitFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails, flags.UseOrgUnit)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

