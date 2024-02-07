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

type removeEntitlementsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	Entitlements          []string
}

func removeEntitlementsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeEntitlements",
		Short: "Remove addon entitlements from a server. Entitlements a server does
 not have are quietly ignored.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeEntitlementsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeEntitlements)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Entitlements", "", "$desc")

	return cmd
}

func removeEntitlements(globalFlags *types.GlobalFlags, flags *removeEntitlementsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.Entitlements)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

