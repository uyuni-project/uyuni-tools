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

type addTrustFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
	TrustOrgId          int
}

func addTrustCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addTrust",
		Short: "Add an organization to the list of trusted organizations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addTrustFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addTrust)
		},
	}

	cmd.Flags().String("OrgId", "", "")
	cmd.Flags().String("TrustOrgId", "", "")

	return cmd
}

func addTrust(globalFlags *types.GlobalFlags, flags *addTrustFlags, cmd *cobra.Command, args []string) error {

res, err := trusts.Trusts(&flags.ConnectionDetails, flags.OrgId, flags.TrustOrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

