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

type schedulePackageUpdateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sids          []int
	EarliestOccurrence          $date
}

func schedulePackageUpdateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedulePackageUpdate",
		Short: "Schedule full package update for several systems.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags schedulePackageUpdateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, schedulePackageUpdate)
		},
	}

	cmd.Flags().String("Sids", "", "$desc")
	cmd.Flags().String("EarliestOccurrence", "", "")

	return cmd
}

func schedulePackageUpdate(globalFlags *types.GlobalFlags, flags *schedulePackageUpdateFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sids, flags.EarliestOccurrence)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

