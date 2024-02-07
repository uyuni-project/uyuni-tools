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

type getGroupFormulaDataFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	GroupId               int
	FormulaName           string
}

func getGroupFormulaDataCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getGroupFormulaData",
		Short: "Get the saved data for the specific formula against specific group",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getGroupFormulaDataFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getGroupFormulaData)
		},
	}

	cmd.Flags().String("GroupId", "", "")
	cmd.Flags().String("FormulaName", "", "")

	return cmd
}

func getGroupFormulaData(globalFlags *types.GlobalFlags, flags *getGroupFormulaDataFlags, cmd *cobra.Command, args []string) error {

	res, err := formula.Formula(&flags.ConnectionDetails, flags.GroupId, flags.FormulaName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
