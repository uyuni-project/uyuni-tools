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

type listSourcePackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listSourcePackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSourcePackages",
		Short: "List all source packages in user's organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSourcePackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSourcePackages)
		},
	}

	return cmd
}

func listSourcePackages(globalFlags *types.GlobalFlags, flags *listSourcePackagesFlags, cmd *cobra.Command, args []string) error {

	res, err := packages.Packages(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
