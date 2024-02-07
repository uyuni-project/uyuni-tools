package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listSystemGroupsForSystemsWithEntitlementFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Entitlement          string
}

func listSystemGroupsForSystemsWithEntitlementCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSystemGroupsForSystemsWithEntitlement",
		Short: "Returns the groups information a system is member of, for all the systems visible to the passed user
 and that are entitled with the passed entitlement.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSystemGroupsForSystemsWithEntitlementFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSystemGroupsForSystemsWithEntitlement)
		},
	}

	cmd.Flags().String("Entitlement", "", "")

	return cmd
}

func listSystemGroupsForSystemsWithEntitlement(globalFlags *types.GlobalFlags, flags *listSystemGroupsForSystemsWithEntitlementFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Entitlement)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

