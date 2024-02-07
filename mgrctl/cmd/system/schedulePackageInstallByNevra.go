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

type schedulePackageInstallByNevraFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sids          []int
	$param.getFlagName()          $param.getType()
	EarliestOccurrence          $date
	$param.getFlagName()          $param.getType()
	AllowModules          bool
	Sid          int
	AllowModules          bool
}

func schedulePackageInstallByNevraCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedulePackageInstallByNevra",
		Short: "Schedule package installation for several systems.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags schedulePackageInstallByNevraFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, schedulePackageInstallByNevra)
		},
	}

	cmd.Flags().String("Sids", "", "$desc")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("EarliestOccurrence", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("AllowModules", "", "Allow this API call, despite modular content being present")
	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("AllowModules", "", "Allow this API call, despite modular content being present")

	return cmd
}

func schedulePackageInstallByNevra(globalFlags *types.GlobalFlags, flags *schedulePackageInstallByNevraFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sids, flags.$param.getFlagName(), flags.EarliestOccurrence, flags.$param.getFlagName(), flags.AllowModules, flags.Sid, flags.AllowModules)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

