package systemgroup

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/systemgroup"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listAssignedFormualsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemGroupName       string
}

func listAssignedFormualsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAssignedFormuals",
		Short: "List all Configuration Channels assigned to a system group",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAssignedFormualsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAssignedFormuals)
		},
	}

	cmd.Flags().String("SystemGroupName", "", "")

	return cmd
}

func listAssignedFormuals(globalFlags *types.GlobalFlags, flags *listAssignedFormualsFlags, cmd *cobra.Command, args []string) error {

	res, err := systemgroup.Systemgroup(&flags.ConnectionDetails, flags.SystemGroupName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
