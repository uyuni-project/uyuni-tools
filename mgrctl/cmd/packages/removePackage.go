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

type removePackageFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Pid          int
}

func removePackageCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removePackage",
		Short: "Remove a package from #product().",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removePackageFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removePackage)
		},
	}

	cmd.Flags().String("Pid", "", "")

	return cmd
}

func removePackage(globalFlags *types.GlobalFlags, flags *removePackageFlags, cmd *cobra.Command, args []string) error {

res, err := packages.Packages(&flags.ConnectionDetails, flags.Pid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

