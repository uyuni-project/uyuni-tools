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

type setFormulasOfServerFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
	Formulas              []string
}

func setFormulasOfServerCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setFormulasOfServer",
		Short: "Set the formulas of a server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setFormulasOfServerFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setFormulasOfServer)
		},
	}

	cmd.Flags().String("Sid", "", "the system ID")
	cmd.Flags().String("Formulas", "", "$desc")

	return cmd
}

func setFormulasOfServer(globalFlags *types.GlobalFlags, flags *setFormulasOfServerFlags, cmd *cobra.Command, args []string) error {

	res, err := formula.Formula(&flags.ConnectionDetails, flags.Sid, flags.Formulas)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
