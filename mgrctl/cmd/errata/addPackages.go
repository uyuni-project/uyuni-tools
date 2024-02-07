package errata

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/errata"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type addPackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	AdvisoryName          string
	PackageIds          []int
}

func addPackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addPackages",
		Short: "Add a set of packages to an erratum with the given advisory name.
 This method will only allow for modification of custom errata created either through the UI or API.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addPackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addPackages)
		},
	}

	cmd.Flags().String("AdvisoryName", "", "")
	cmd.Flags().String("PackageIds", "", "$desc")

	return cmd
}

func addPackages(globalFlags *types.GlobalFlags, flags *addPackagesFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.AdvisoryName, flags.PackageIds)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

