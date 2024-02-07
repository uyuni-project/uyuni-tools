package proxy

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/proxy"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type createMonitoringScoutFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Clientcert          string
}

func createMonitoringScoutCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createMonitoringScout",
		Short: "Create Monitoring Scout for proxy.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createMonitoringScoutFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createMonitoringScout)
		},
	}

	cmd.Flags().String("Clientcert", "", "client certificate file")

	return cmd
}

func createMonitoringScout(globalFlags *types.GlobalFlags, flags *createMonitoringScoutFlags, cmd *cobra.Command, args []string) error {

res, err := proxy.Proxy(&flags.ConnectionDetails, flags.Clientcert)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

