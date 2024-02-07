package formula

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/formula"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getSystemFormulaDataFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	FormulaName          string
}

func getSystemFormulaDataCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getSystemFormulaData",
		Short: "Get the saved data for the specific formula against specific server",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getSystemFormulaDataFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getSystemFormulaData)
		},
	}

	cmd.Flags().String("Sid", "", "the system ID")
	cmd.Flags().String("FormulaName", "", "")

	return cmd
}

func getSystemFormulaData(globalFlags *types.GlobalFlags, flags *getSystemFormulaDataFlags, cmd *cobra.Command, args []string) error {

res, err := formula.Formula(&flags.ConnectionDetails, flags.Sid, flags.FormulaName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

