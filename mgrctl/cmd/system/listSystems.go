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

type listSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSystems",
		Short: "Returns a list of all servers visible to the user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSystems)
		},
	}

	return cmd
}

func listSystems(globalFlags *types.GlobalFlags, flags *listSystemsFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
