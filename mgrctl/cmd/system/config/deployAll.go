package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/config"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type deployAllFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Date          $type
}

func deployAllCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deployAll",
		Short: "Schedules a deploy action for all the configuration files
 on the given list of systems.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deployAllFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deployAll)
		},
	}

	cmd.Flags().String("Date", "", "Earliest date for the deploy action.")

	return cmd
}

func deployAll(globalFlags *types.GlobalFlags, flags *deployAllFlags, cmd *cobra.Command, args []string) error {

res, err := config.Config(&flags.ConnectionDetails, flags.Date)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

