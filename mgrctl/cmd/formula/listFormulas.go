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

type listFormulasFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listFormulasCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listFormulas",
		Short: "Return the list of formulas currently installed.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listFormulasFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listFormulas)
		},
	}

	return cmd
}

func listFormulas(globalFlags *types.GlobalFlags, flags *listFormulasFlags, cmd *cobra.Command, args []string) error {

	res, err := formula.Formula(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
