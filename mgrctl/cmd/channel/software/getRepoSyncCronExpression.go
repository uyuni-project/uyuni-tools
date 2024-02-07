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

type getRepoSyncCronExpressionFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func getRepoSyncCronExpressionCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getRepoSyncCronExpression",
		Short: "Returns repo synchronization cron expression",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getRepoSyncCronExpressionFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getRepoSyncCronExpression)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "channel label")

	return cmd
}

func getRepoSyncCronExpression(globalFlags *types.GlobalFlags, flags *getRepoSyncCronExpressionFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
