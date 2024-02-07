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

type isNvreInstalledFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
	Name                  string
	Version               string
	Release               string
	Version               string
	Epoch                 string
}

func isNvreInstalledCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "isNvreInstalled",
		Short: "Check if the package with the given NVRE is installed on given system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags isNvreInstalledFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, isNvreInstalled)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Name", "", "Package name.")
	cmd.Flags().String("Version", "", "Package version.")
	cmd.Flags().String("Release", "", "Package release.")
	cmd.Flags().String("Version", "", "Package version.")
	cmd.Flags().String("Epoch", "", "Package epoch.")

	return cmd
}

func isNvreInstalled(globalFlags *types.GlobalFlags, flags *isNvreInstalledFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.Name, flags.Version, flags.Release, flags.Version, flags.Epoch)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
