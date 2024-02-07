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

type getDefaultOrgFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func getDefaultOrgCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getDefaultOrg",
		Short: "Get the default org that users should be added in if orgunit from
 IPA server isn't found or is disabled. Can only be called by a #product() Administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getDefaultOrgFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getDefaultOrg)
		},
	}


	return cmd
}

func getDefaultOrg(globalFlags *types.GlobalFlags, flags *getDefaultOrgFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

