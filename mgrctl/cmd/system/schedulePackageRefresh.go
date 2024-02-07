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

type schedulePackageRefreshFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	EarliestOccurrence          $date
}

func schedulePackageRefreshCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedulePackageRefresh",
		Short: "Schedule a package list refresh for a system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags schedulePackageRefreshFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, schedulePackageRefresh)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("EarliestOccurrence", "", "")

	return cmd
}

func schedulePackageRefresh(globalFlags *types.GlobalFlags, flags *schedulePackageRefreshFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.EarliestOccurrence)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

