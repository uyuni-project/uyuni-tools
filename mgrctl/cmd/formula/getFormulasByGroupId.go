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

type getFormulasByGroupIdFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemGroupId          int
}

func getFormulasByGroupIdCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getFormulasByGroupId",
		Short: "Return the list of formulas a server group has.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getFormulasByGroupIdFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getFormulasByGroupId)
		},
	}

	cmd.Flags().String("SystemGroupId", "", "")

	return cmd
}

func getFormulasByGroupId(globalFlags *types.GlobalFlags, flags *getFormulasByGroupIdFlags, cmd *cobra.Command, args []string) error {

res, err := formula.Formula(&flags.ConnectionDetails, flags.SystemGroupId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

