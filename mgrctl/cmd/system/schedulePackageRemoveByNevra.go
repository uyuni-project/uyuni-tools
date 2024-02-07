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

type schedulePackageRemoveByNevraFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sids          []int
	$param.getFlagName()          $param.getType()
	EarliestOccurrence          $date
	AllowModules          bool
	Sid          int
}

func schedulePackageRemoveByNevraCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedulePackageRemoveByNevra",
		Short: "Schedule package removal for several systems.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags schedulePackageRemoveByNevraFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, schedulePackageRemoveByNevra)
		},
	}

	cmd.Flags().String("Sids", "", "$desc")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("EarliestOccurrence", "", "")
	cmd.Flags().String("AllowModules", "", "Allow this API call, despite modular content being present")
	cmd.Flags().String("Sid", "", "")

	return cmd
}

func schedulePackageRemoveByNevra(globalFlags *types.GlobalFlags, flags *schedulePackageRemoveByNevraFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sids, flags.$param.getFlagName(), flags.EarliestOccurrence, flags.AllowModules, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

