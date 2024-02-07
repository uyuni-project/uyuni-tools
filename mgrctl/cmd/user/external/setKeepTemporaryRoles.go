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

type setKeepTemporaryRolesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KeepRoles          bool
}

func setKeepTemporaryRolesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setKeepTemporaryRoles",
		Short: "Set whether we should keeps roles assigned to users because of
 their IPA groups even after they log in through a non-IPA method. Can only be
 called by a #product() Administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setKeepTemporaryRolesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setKeepTemporaryRoles)
		},
	}

	cmd.Flags().String("KeepRoles", "", "True if we should keep roles after users log in through non-IPA method, false otherwise.")

	return cmd
}

func setKeepTemporaryRoles(globalFlags *types.GlobalFlags, flags *setKeepTemporaryRolesFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails, flags.KeepRoles)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

