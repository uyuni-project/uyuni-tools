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

type listLatestUpgradablePackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func listLatestUpgradablePackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listLatestUpgradablePackages",
		Short: "Get the list of latest upgradable packages for a given system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listLatestUpgradablePackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listLatestUpgradablePackages)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listLatestUpgradablePackages(globalFlags *types.GlobalFlags, flags *listLatestUpgradablePackagesFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
