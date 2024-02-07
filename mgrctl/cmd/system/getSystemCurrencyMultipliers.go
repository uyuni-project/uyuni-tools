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

type getSystemCurrencyMultipliersFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func getSystemCurrencyMultipliersCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getSystemCurrencyMultipliers",
		Short: "Get the System Currency score multipliers",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getSystemCurrencyMultipliersFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getSystemCurrencyMultipliers)
		},
	}


	return cmd
}

func getSystemCurrencyMultipliers(globalFlags *types.GlobalFlags, flags *getSystemCurrencyMultipliersFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

