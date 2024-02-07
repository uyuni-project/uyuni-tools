package actionchain

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/actionchain"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type addPackageInstallFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	PackageIds          []int
	ChainLabel          string
}

func addPackageInstallCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addPackageInstall",
		Short: "Adds package installation action to an Action Chain.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addPackageInstallFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addPackageInstall)
		},
	}

	cmd.Flags().String("Sid", "", "System ID")
	cmd.Flags().String("PackageIds", "", "$desc")
	cmd.Flags().String("ChainLabel", "", "")

	return cmd
}

func addPackageInstall(globalFlags *types.GlobalFlags, flags *addPackageInstallFlags, cmd *cobra.Command, args []string) error {

res, err := actionchain.Actionchain(&flags.ConnectionDetails, flags.Sid, flags.PackageIds, flags.ChainLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

