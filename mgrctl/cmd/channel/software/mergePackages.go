package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type mergePackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	MergeFromLabel        string
	MergeToLabel          string
	AlignModules          bool
}

func mergePackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mergePackages",
		Short: "Merges all packages from one channel into another",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags mergePackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, mergePackages)
		},
	}

	cmd.Flags().String("MergeFromLabel", "", "the label of the          channel to pull packages from")
	cmd.Flags().String("MergeToLabel", "", "the label to push the              packages into")
	cmd.Flags().String("AlignModules", "", "align modular data of the target channel              to the source channel (RHEL8 and higher)")

	return cmd
}

func mergePackages(globalFlags *types.GlobalFlags, flags *mergePackagesFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.MergeFromLabel, flags.MergeToLabel, flags.AlignModules)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
