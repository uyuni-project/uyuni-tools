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

type addEntitlementsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	Entitlements          []string
}

func addEntitlementsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addEntitlements",
		Short: "Add entitlements to a server. Entitlements a server already has
 are quietly ignored.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addEntitlementsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addEntitlements)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Entitlements", "", "$desc")

	return cmd
}

func addEntitlements(globalFlags *types.GlobalFlags, flags *addEntitlementsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.Entitlements)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

