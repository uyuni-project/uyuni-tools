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

type listActiveSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listActiveSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listActiveSystems",
		Short: "Returns a list of active servers visible to the user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listActiveSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listActiveSystems)
		},
	}


	return cmd
}

func listActiveSystems(globalFlags *types.GlobalFlags, flags *listActiveSystemsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

