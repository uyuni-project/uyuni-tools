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

type listPackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	AdvisoryName          string
}

func listPackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listPackages",
		Short: "Returns a list of the packages affected by the errata with the given advisory name.
 For those errata that are present in both vendor and user organizations under the same advisory name,
 this method retrieves the packages of both of them.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listPackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listPackages)
		},
	}

	cmd.Flags().String("AdvisoryName", "", "")

	return cmd
}

func listPackages(globalFlags *types.GlobalFlags, flags *listPackagesFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.AdvisoryName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

