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

type listLatestAvailablePackageFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sids                  []int
	PackageName           string
}

func listLatestAvailablePackageCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listLatestAvailablePackage",
		Short: "Get the latest available version of a package for each system",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listLatestAvailablePackageFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listLatestAvailablePackage)
		},
	}

	cmd.Flags().String("Sids", "", "$desc")
	cmd.Flags().String("PackageName", "", "")

	return cmd
}

func listLatestAvailablePackage(globalFlags *types.GlobalFlags, flags *listLatestAvailablePackageFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sids, flags.PackageName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
