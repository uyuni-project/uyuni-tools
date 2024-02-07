package packages

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/packages"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listDependenciesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Pid          int
}

func listDependenciesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listDependencies",
		Short: "List the dependencies for a package.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listDependenciesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listDependencies)
		},
	}

	cmd.Flags().String("Pid", "", "")

	return cmd
}

func listDependencies(globalFlags *types.GlobalFlags, flags *listDependenciesFlags, cmd *cobra.Command, args []string) error {

res, err := packages.Packages(&flags.ConnectionDetails, flags.Pid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

