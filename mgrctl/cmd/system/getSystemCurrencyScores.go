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

type getSystemCurrencyScoresFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func getSystemCurrencyScoresCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getSystemCurrencyScores",
		Short: "Get the System Currency scores for all servers the user has access to",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getSystemCurrencyScoresFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getSystemCurrencyScores)
		},
	}


	return cmd
}

func getSystemCurrencyScores(globalFlags *types.GlobalFlags, flags *getSystemCurrencyScoresFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

