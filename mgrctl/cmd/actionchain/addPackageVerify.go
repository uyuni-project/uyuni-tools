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

type addPackageVerifyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	PackageIds          []int
	ChainLabel          string
}

func addPackageVerifyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addPackageVerify",
		Short: "Adds an action to verify installed packages on the system to an Action
 Chain.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addPackageVerifyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addPackageVerify)
		},
	}

	cmd.Flags().String("Sid", "", "System ID")
	cmd.Flags().String("PackageIds", "", "$desc")
	cmd.Flags().String("ChainLabel", "", "Label of the chain")

	return cmd
}

func addPackageVerify(globalFlags *types.GlobalFlags, flags *addPackageVerifyFlags, cmd *cobra.Command, args []string) error {

res, err := actionchain.Actionchain(&flags.ConnectionDetails, flags.Sid, flags.PackageIds, flags.ChainLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

