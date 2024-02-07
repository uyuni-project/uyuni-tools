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

type scheduleSyncPackagesWithSystemFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	TargetServerId          int
	SourceServerId          int
	$param.getFlagName()          $param.getType()
	EarliestOccurrence          $date
}

func scheduleSyncPackagesWithSystemCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleSyncPackagesWithSystem",
		Short: "Sync packages from a source system to a target.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleSyncPackagesWithSystemFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleSyncPackagesWithSystem)
		},
	}

	cmd.Flags().String("TargetServerId", "", "Target system to apply package                  changes to.")
	cmd.Flags().String("SourceServerId", "", "Source system to retrieve                  package state from.")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("EarliestOccurrence", "", "Date to schedule action for")

	return cmd
}

func scheduleSyncPackagesWithSystem(globalFlags *types.GlobalFlags, flags *scheduleSyncPackagesWithSystemFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.TargetServerId, flags.SourceServerId, flags.$param.getFlagName(), flags.EarliestOccurrence)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

