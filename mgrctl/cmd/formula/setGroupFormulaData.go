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

type setGroupFormulaDataFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	GroupId          int
	FormulaName          string
	Content          struct
}

func setGroupFormulaDataCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setGroupFormulaData",
		Short: "Set the formula form for the specified group.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setGroupFormulaDataFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setGroupFormulaData)
		},
	}

	cmd.Flags().String("GroupId", "", "")
	cmd.Flags().String("FormulaName", "", "")
	cmd.Flags().String("Content", "", "struct containing the values for each field in the form")

	return cmd
}

func setGroupFormulaData(globalFlags *types.GlobalFlags, flags *setGroupFormulaDataFlags, cmd *cobra.Command, args []string) error {

res, err := formula.Formula(&flags.ConnectionDetails, flags.GroupId, flags.FormulaName, flags.Content)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

