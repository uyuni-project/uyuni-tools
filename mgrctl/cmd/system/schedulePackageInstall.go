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

type schedulePackageInstallFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sids          []int
	PackageIds          []int
	EarliestOccurrence          $date
	AllowModules          bool
	Sid          int
}

func schedulePackageInstallCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedulePackageInstall",
		Short: "Schedule package installation for several systems.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags schedulePackageInstallFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, schedulePackageInstall)
		},
	}

	cmd.Flags().String("Sids", "", "$desc")
	cmd.Flags().String("PackageIds", "", "$desc")
	cmd.Flags().String("EarliestOccurrence", "", "")
	cmd.Flags().String("AllowModules", "", "Allow this API call, despite modular content being present")
	cmd.Flags().String("Sid", "", "")

	return cmd
}

func schedulePackageInstall(globalFlags *types.GlobalFlags, flags *schedulePackageInstallFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sids, flags.PackageIds, flags.EarliestOccurrence, flags.AllowModules, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

