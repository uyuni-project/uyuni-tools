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

type getEntitlementsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func getEntitlementsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getEntitlements",
		Short: "Gets the entitlements for a given server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getEntitlementsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getEntitlements)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getEntitlements(globalFlags *types.GlobalFlags, flags *getEntitlementsFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
