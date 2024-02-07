package trusts

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/org/trusts"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type removeTrustFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
	TrustOrgId          int
}

func removeTrustCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeTrust",
		Short: "Remove an organization to the list of trusted organizations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeTrustFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeTrust)
		},
	}

	cmd.Flags().String("OrgId", "", "")
	cmd.Flags().String("TrustOrgId", "", "")

	return cmd
}

func removeTrust(globalFlags *types.GlobalFlags, flags *removeTrustFlags, cmd *cobra.Command, args []string) error {

res, err := trusts.Trusts(&flags.ConnectionDetails, flags.OrgId, flags.TrustOrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

