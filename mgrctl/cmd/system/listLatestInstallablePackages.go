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

type listLatestInstallablePackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func listLatestInstallablePackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listLatestInstallablePackages",
		Short: "Get the list of latest installable packages for a given system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listLatestInstallablePackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listLatestInstallablePackages)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listLatestInstallablePackages(globalFlags *types.GlobalFlags, flags *listLatestInstallablePackagesFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

