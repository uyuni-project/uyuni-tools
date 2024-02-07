package formula

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "formula",
		Short: "Provides methods to access and modify formulas.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(getCombinedFormulasByServerIdCommand(globalFlags))
	cmd.AddCommand(getFormulasByServerIdCommand(globalFlags))
	cmd.AddCommand(setFormulasOfGroupCommand(globalFlags))
	cmd.AddCommand(setGroupFormulaDataCommand(globalFlags))
	cmd.AddCommand(setFormulasOfServerCommand(globalFlags))
	cmd.AddCommand(getFormulasByGroupIdCommand(globalFlags))
	cmd.AddCommand(setSystemFormulaDataCommand(globalFlags))
	cmd.AddCommand(getSystemFormulaDataCommand(globalFlags))
	cmd.AddCommand(getCombinedFormulaDataByServerIdsCommand(globalFlags))
	cmd.AddCommand(listFormulasCommand(globalFlags))
	cmd.AddCommand(getGroupFormulaDataCommand(globalFlags))

	return cmd
}
