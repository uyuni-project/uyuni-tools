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

type scheduleDistUpgradeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	Channels          []string
	DryRun          bool
	EarliestOccurrence          $date
	AllowVendorChange          bool
}

func scheduleDistUpgradeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleDistUpgrade",
		Short: "Schedule a dist upgrade for a system. This call takes a list of channel
 labels that the system will be subscribed to before performing the dist upgrade.
 Note: You can seriously damage your system with this call, use it only if you really
 know what you are doing! Make sure that the list of channel labels is complete and in
 any case do a dry run before scheduling an actual dist upgrade.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleDistUpgradeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleDistUpgrade)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Channels", "", "$desc")
	cmd.Flags().String("DryRun", "", "")
	cmd.Flags().String("EarliestOccurrence", "", "")
	cmd.Flags().String("AllowVendorChange", "", "")

	return cmd
}

func scheduleDistUpgrade(globalFlags *types.GlobalFlags, flags *scheduleDistUpgradeFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.Channels, flags.DryRun, flags.EarliestOccurrence, flags.AllowVendorChange)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

