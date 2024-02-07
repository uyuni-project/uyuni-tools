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

type upgradeEntitlementFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	EntitlementLevel          string
}

func upgradeEntitlementCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgradeEntitlement",
		Short: "Adds an entitlement to a given server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags upgradeEntitlementFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, upgradeEntitlement)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("EntitlementLevel", "", "One of:          'enterprise_entitled' or 'virtualization_host'.")

	return cmd
}

func upgradeEntitlement(globalFlags *types.GlobalFlags, flags *upgradeEntitlementFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.EntitlementLevel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

