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

type getCombinedFormulasByServerIdFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func getCombinedFormulasByServerIdCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getCombinedFormulasByServerId",
		Short: "Return the list of formulas a server and all his groups have.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getCombinedFormulasByServerIdFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getCombinedFormulasByServerId)
		},
	}

	cmd.Flags().String("Sid", "", "the system ID")

	return cmd
}

func getCombinedFormulasByServerId(globalFlags *types.GlobalFlags, flags *getCombinedFormulasByServerIdFlags, cmd *cobra.Command, args []string) error {

res, err := formula.Formula(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

