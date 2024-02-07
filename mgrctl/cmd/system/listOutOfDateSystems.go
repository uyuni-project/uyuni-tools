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

type listOutOfDateSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listOutOfDateSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listOutOfDateSystems",
		Short: "Returns list of systems needing package updates.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listOutOfDateSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listOutOfDateSystems)
		},
	}


	return cmd
}

func listOutOfDateSystems(globalFlags *types.GlobalFlags, flags *listOutOfDateSystemsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

