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

type listSystemsWithExtraPackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listSystemsWithExtraPackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSystemsWithExtraPackages",
		Short: "List systems with extra packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSystemsWithExtraPackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSystemsWithExtraPackages)
		},
	}

	return cmd
}

func listSystemsWithExtraPackages(globalFlags *types.GlobalFlags, flags *listSystemsWithExtraPackagesFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
