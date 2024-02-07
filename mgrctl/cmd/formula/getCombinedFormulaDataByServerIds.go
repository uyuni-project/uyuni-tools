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

type getCombinedFormulaDataByServerIdsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	FormulaName           string
	Sids                  []int
}

func getCombinedFormulaDataByServerIdsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getCombinedFormulaDataByServerIds",
		Short: "Return the list of formulas a server and all his groups have.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getCombinedFormulaDataByServerIdsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getCombinedFormulaDataByServerIds)
		},
	}

	cmd.Flags().String("FormulaName", "", "")
	cmd.Flags().String("Sids", "", "$desc")

	return cmd
}

func getCombinedFormulaDataByServerIds(globalFlags *types.GlobalFlags, flags *getCombinedFormulaDataByServerIdsFlags, cmd *cobra.Command, args []string) error {

	res, err := formula.Formula(&flags.ConnectionDetails, flags.FormulaName, flags.Sids)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
