package configchannel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/configchannel"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type deployAllSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	Date          $date
	FilePath          string
	Date          $date
}

func deployAllSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deployAllSystems",
		Short: "Schedule an immediate configuration deployment for all systems
    subscribed to a particular configuration channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deployAllSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deployAllSystems)
		},
	}

	cmd.Flags().String("Label", "", "the configuration channel's label")
	cmd.Flags().String("Date", "", "the date to schedule the action")
	cmd.Flags().String("FilePath", "", "the configuration file path")
	cmd.Flags().String("Date", "", "the date to schedule the action")

	return cmd
}

func deployAllSystems(globalFlags *types.GlobalFlags, flags *deployAllSystemsFlags, cmd *cobra.Command, args []string) error {

res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label, flags.Date, flags.FilePath, flags.Date)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

