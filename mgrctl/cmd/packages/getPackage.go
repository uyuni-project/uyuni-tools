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

type getPackageFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Pid          int
}

func getPackageCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getPackage",
		Short: "Retrieve the package file associated with a package.
 (Consider using #getPackageUrlpackages.getPackageUrl
 for larger files.)",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getPackageFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getPackage)
		},
	}

	cmd.Flags().String("Pid", "", "")

	return cmd
}

func getPackage(globalFlags *types.GlobalFlags, flags *getPackageFlags, cmd *cobra.Command, args []string) error {

res, err := packages.Packages(&flags.ConnectionDetails, flags.Pid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

