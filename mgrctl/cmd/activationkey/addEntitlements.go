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

type addEntitlementsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Key                   string
}

func addEntitlementsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addEntitlements",
		Short: "Add add-on System Types to an activation key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addEntitlementsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addEntitlements)
		},
	}

	cmd.Flags().String("Key", "", "")

	return cmd
}

func addEntitlements(globalFlags *types.GlobalFlags, flags *addEntitlementsFlags, cmd *cobra.Command, args []string) error {

	res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Key)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
