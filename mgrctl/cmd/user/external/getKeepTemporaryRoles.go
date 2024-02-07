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

type getKeepTemporaryRolesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func getKeepTemporaryRolesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getKeepTemporaryRoles",
		Short: "Get whether we should keeps roles assigned to users because of
 their IPA groups even after they log in through a non-IPA method. Can only be
 called by a #product() Administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getKeepTemporaryRolesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getKeepTemporaryRoles)
		},
	}


	return cmd
}

func getKeepTemporaryRoles(globalFlags *types.GlobalFlags, flags *getKeepTemporaryRolesFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

