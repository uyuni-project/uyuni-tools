package actionchain

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/actionchain"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type scheduleChainFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChainLabel          string
	Date          $date
}

func scheduleChainCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleChain",
		Short: "Schedule the Action Chain so that its actions will actually occur.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleChainFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleChain)
		},
	}

	cmd.Flags().String("ChainLabel", "", "Label of the chain")
	cmd.Flags().String("Date", "", "Earliest date")

	return cmd
}

func scheduleChain(globalFlags *types.GlobalFlags, flags *scheduleChainFlags, cmd *cobra.Command, args []string) error {

res, err := actionchain.Actionchain(&flags.ConnectionDetails, flags.ChainLabel, flags.Date)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

