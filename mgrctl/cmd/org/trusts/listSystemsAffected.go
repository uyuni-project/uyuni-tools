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

type listSystemsAffectedFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
	TrustOrgId          string
}

func listSystemsAffectedCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSystemsAffected",
		Short: "Get a list of systems within the  trusted organization
   that would be affected if the trust relationship was removed.
   This basically lists systems that are sharing at least (1) package.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSystemsAffectedFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSystemsAffected)
		},
	}

	cmd.Flags().String("OrgId", "", "")
	cmd.Flags().String("TrustOrgId", "", "")

	return cmd
}

func listSystemsAffected(globalFlags *types.GlobalFlags, flags *listSystemsAffectedFlags, cmd *cobra.Command, args []string) error {

res, err := trusts.Trusts(&flags.ConnectionDetails, flags.OrgId, flags.TrustOrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

