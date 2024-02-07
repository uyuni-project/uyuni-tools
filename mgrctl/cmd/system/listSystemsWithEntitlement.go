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

type listSystemsWithEntitlementFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	EntitlementName          string
}

func listSystemsWithEntitlementCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSystemsWithEntitlement",
		Short: "Lists the systems that have the given entitlement",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSystemsWithEntitlementFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSystemsWithEntitlement)
		},
	}

	cmd.Flags().String("EntitlementName", "", "the entitlement name")

	return cmd
}

func listSystemsWithEntitlement(globalFlags *types.GlobalFlags, flags *listSystemsWithEntitlementFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.EntitlementName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

