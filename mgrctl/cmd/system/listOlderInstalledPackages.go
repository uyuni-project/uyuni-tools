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

type listOlderInstalledPackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	Name          string
	Version          string
	Release          string
	Epoch          string
}

func listOlderInstalledPackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listOlderInstalledPackages",
		Short: "Given a package name, version, release, and epoch, returns
 the list of packages installed on the system with the same name that are
 older.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listOlderInstalledPackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listOlderInstalledPackages)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Name", "", "Package name.")
	cmd.Flags().String("Version", "", "Package version.")
	cmd.Flags().String("Release", "", "Package release.")
	cmd.Flags().String("Epoch", "", "Package epoch.")

	return cmd
}

func listOlderInstalledPackages(globalFlags *types.GlobalFlags, flags *listOlderInstalledPackagesFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.Name, flags.Version, flags.Release, flags.Epoch)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

