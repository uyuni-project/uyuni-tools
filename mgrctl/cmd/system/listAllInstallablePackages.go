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

type listAllInstallablePackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func listAllInstallablePackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAllInstallablePackages",
		Short: "Get the list of all installable packages for a given system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAllInstallablePackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAllInstallablePackages)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listAllInstallablePackages(globalFlags *types.GlobalFlags, flags *listAllInstallablePackagesFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
