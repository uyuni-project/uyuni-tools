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

type listActiveSystemsInGroupFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemGroupName          string
}

func listActiveSystemsInGroupCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listActiveSystemsInGroup",
		Short: "Lists active systems within a server group",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listActiveSystemsInGroupFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listActiveSystemsInGroup)
		},
	}

	cmd.Flags().String("SystemGroupName", "", "")

	return cmd
}

func listActiveSystemsInGroup(globalFlags *types.GlobalFlags, flags *listActiveSystemsInGroupFlags, cmd *cobra.Command, args []string) error {

res, err := systemgroup.Systemgroup(&flags.ConnectionDetails, flags.SystemGroupName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

