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

type listInstalledPackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func listInstalledPackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listInstalledPackages",
		Short: "List the installed packages for a given system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listInstalledPackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listInstalledPackages)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listInstalledPackages(globalFlags *types.GlobalFlags, flags *listInstalledPackagesFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

