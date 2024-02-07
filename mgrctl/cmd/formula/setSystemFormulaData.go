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

type setSystemFormulaDataFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemId          int
	FormulaName          string
	Content          struct
}

func setSystemFormulaDataCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setSystemFormulaData",
		Short: "Set the formula form for the specified server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setSystemFormulaDataFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setSystemFormulaData)
		},
	}

	cmd.Flags().String("SystemId", "", "")
	cmd.Flags().String("FormulaName", "", "")
	cmd.Flags().String("Content", "", "struct content with the values for each field in the form")

	return cmd
}

func setSystemFormulaData(globalFlags *types.GlobalFlags, flags *setSystemFormulaDataFlags, cmd *cobra.Command, args []string) error {

res, err := formula.Formula(&flags.ConnectionDetails, flags.SystemId, flags.FormulaName, flags.Content)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

