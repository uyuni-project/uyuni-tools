package activationkey

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/activationkey"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type removeEntitlementsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Key          string
}

func removeEntitlementsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeEntitlements",
		Short: "Remove entitlements (by label) from an activation key.
 Currently only virtualization_host add-on entitlement is permitted.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeEntitlementsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeEntitlements)
		},
	}

	cmd.Flags().String("Key", "", "")

	return cmd
}

func removeEntitlements(globalFlags *types.GlobalFlags, flags *removeEntitlementsFlags, cmd *cobra.Command, args []string) error {

res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Key)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

