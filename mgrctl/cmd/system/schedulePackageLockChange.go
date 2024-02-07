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

type schedulePackageLockChangeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	PkgIdsToLock          []int
	PkgIdsToUnlock          []int
	EarliestOccurrence          $date
}

func schedulePackageLockChangeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedulePackageLockChange",
		Short: "Schedule package lock for a system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags schedulePackageLockChangeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, schedulePackageLockChange)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("PkgIdsToLock", "", "$desc")
	cmd.Flags().String("PkgIdsToUnlock", "", "$desc")
	cmd.Flags().String("EarliestOccurrence", "", "")

	return cmd
}

func schedulePackageLockChange(globalFlags *types.GlobalFlags, flags *schedulePackageLockChangeFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.PkgIdsToLock, flags.PkgIdsToUnlock, flags.EarliestOccurrence)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

