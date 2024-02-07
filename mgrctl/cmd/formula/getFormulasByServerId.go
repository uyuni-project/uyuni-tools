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

type getFormulasByServerIdFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func getFormulasByServerIdCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getFormulasByServerId",
		Short: "Return the list of formulas directly applied to a server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getFormulasByServerIdFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getFormulasByServerId)
		},
	}

	cmd.Flags().String("Sid", "", "the system ID")

	return cmd
}

func getFormulasByServerId(globalFlags *types.GlobalFlags, flags *getFormulasByServerIdFlags, cmd *cobra.Command, args []string) error {

	res, err := formula.Formula(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
