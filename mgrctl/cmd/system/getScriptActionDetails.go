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

type getScriptActionDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ActionId          int
}

func getScriptActionDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getScriptActionDetails",
		Short: "Returns script details for script run actions",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getScriptActionDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getScriptActionDetails)
		},
	}

	cmd.Flags().String("ActionId", "", "ID of the script run action.")

	return cmd
}

func getScriptActionDetails(globalFlags *types.GlobalFlags, flags *getScriptActionDetailsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.ActionId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

