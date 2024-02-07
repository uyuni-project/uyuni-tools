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

type scheduleProductMigrationFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	BaseChannelLabel          string
	OptionalChildChannels          []string
	DryRun          bool
	EarliestOccurrence          $date
	AllowVendorChange          bool
	TargetIdent          string
	EarliestOccurrence          $date
	TargetIdent          string
	RemoveProductsWithNoSuccessorAfterMigration          bool
}

func scheduleProductMigrationCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleProductMigration",
		Short: "Schedule a Product migration for a system. This call is the
 recommended and supported way of migrating a system to the next Service Pack. It will
 automatically find all mandatory product channels below a given target base channel
 and subscribe the system accordingly. Any additional optional channels can be
 subscribed by providing their labels.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleProductMigrationFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleProductMigration)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("BaseChannelLabel", "", "")
	cmd.Flags().String("OptionalChildChannels", "", "$desc")
	cmd.Flags().String("DryRun", "", "")
	cmd.Flags().String("EarliestOccurrence", "", "")
	cmd.Flags().String("AllowVendorChange", "", "")
	cmd.Flags().String("TargetIdent", "", "Identifier for the selected migration target. Use listMigrationTargets to list the identifiers")
	cmd.Flags().String("EarliestOccurrence", "", "")
	cmd.Flags().String("TargetIdent", "", "Identifier for the selected migration target - User listMigrationTargets to list the identifiers ")
	cmd.Flags().String("RemoveProductsWithNoSuccessorAfterMigration", "", "set to remove products which have no successors. This flag will only have effect if targetIdent will also be specified")

	return cmd
}

func scheduleProductMigration(globalFlags *types.GlobalFlags, flags *scheduleProductMigrationFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.BaseChannelLabel, flags.OptionalChildChannels, flags.DryRun, flags.EarliestOccurrence, flags.AllowVendorChange, flags.TargetIdent, flags.EarliestOccurrence, flags.TargetIdent, flags.RemoveProductsWithNoSuccessorAfterMigration)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

