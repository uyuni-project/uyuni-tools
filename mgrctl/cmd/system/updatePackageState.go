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

type updatePackageStateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	PackageName          string
	State          int
	VersionConstraint          int
}

func updatePackageStateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updatePackageState",
		Short: "Update the package state of a given system
                          (High state would be needed to actually install/remove the package)",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updatePackageStateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, updatePackageState)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("PackageName", "", "Name of the package")
	cmd.Flags().String("State", "", "0 = installed, 1 = removed, 2 = unmanaged ")
	cmd.Flags().String("VersionConstraint", "", "0 = latest, 1 = any ")

	return cmd
}

func updatePackageState(globalFlags *types.GlobalFlags, flags *updatePackageStateFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.PackageName, flags.State, flags.VersionConstraint)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

