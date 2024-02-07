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

type setDefaultOrgFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
}

func setDefaultOrgCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setDefaultOrg",
		Short: "Set the default org that users should be added in if orgunit from
 IPA server isn't found or is disabled. Can only be called by a #product() Administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setDefaultOrgFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setDefaultOrg)
		},
	}

	cmd.Flags().String("OrgId", "", "ID of the organization to set as the default org. 0 if there should not be a default organization.")

	return cmd
}

func setDefaultOrg(globalFlags *types.GlobalFlags, flags *setDefaultOrgFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails, flags.OrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

