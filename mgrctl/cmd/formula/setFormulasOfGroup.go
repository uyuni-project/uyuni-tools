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

type setFormulasOfGroupFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemGroupId          int
	Formulas          []string
}

func setFormulasOfGroupCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setFormulasOfGroup",
		Short: "Set the formulas of a server group.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setFormulasOfGroupFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setFormulasOfGroup)
		},
	}

	cmd.Flags().String("SystemGroupId", "", "")
	cmd.Flags().String("Formulas", "", "$desc")

	return cmd
}

func setFormulasOfGroup(globalFlags *types.GlobalFlags, flags *setFormulasOfGroupFlags, cmd *cobra.Command, args []string) error {

res, err := formula.Formula(&flags.ConnectionDetails, flags.SystemGroupId, flags.Formulas)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

