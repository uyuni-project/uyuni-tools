package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getScriptResultsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ActionId          int
}

func getScriptResultsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getScriptResults",
		Short: "Fetch results from a script execution. Returns an empty array if no
 results are yet available.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getScriptResultsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getScriptResults)
		},
	}

	cmd.Flags().String("ActionId", "", "ID of the script run action.")

	return cmd
}

func getScriptResults(globalFlags *types.GlobalFlags, flags *getScriptResultsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.ActionId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

