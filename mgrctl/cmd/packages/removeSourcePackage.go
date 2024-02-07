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

type removeSourcePackageFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Psid          int
}

func removeSourcePackageCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeSourcePackage",
		Short: "Remove a source package.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeSourcePackageFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeSourcePackage)
		},
	}

	cmd.Flags().String("Psid", "", "package source ID")

	return cmd
}

func removeSourcePackage(globalFlags *types.GlobalFlags, flags *removeSourcePackageFlags, cmd *cobra.Command, args []string) error {

res, err := packages.Packages(&flags.ConnectionDetails, flags.Psid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

