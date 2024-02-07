package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type syncRepoFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabels          []string
	ChannelLabel          string
	$param.getFlagName()          $param.getType()
	CronExpr          string
}

func syncRepoCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "syncRepo",
		Short: "Trigger immediate repo synchronization",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags syncRepoFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, syncRepo)
		},
	}

	cmd.Flags().String("ChannelLabels", "", "$desc")
	cmd.Flags().String("ChannelLabel", "", "channel label")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("CronExpr", "", "cron expression, if empty all periodic schedules will be disabled")

	return cmd
}

func syncRepo(globalFlags *types.GlobalFlags, flags *syncRepoFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabels, flags.ChannelLabel, flags.$param.getFlagName(), flags.CronExpr)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

