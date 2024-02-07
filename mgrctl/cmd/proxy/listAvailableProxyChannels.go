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

type listAvailableProxyChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Clientcert          string
}

func listAvailableProxyChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAvailableProxyChannels",
		Short: "List available version of proxy channel for system
 identified by the given client certificate i.e. systemid file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAvailableProxyChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAvailableProxyChannels)
		},
	}

	cmd.Flags().String("Clientcert", "", "client certificate file")

	return cmd
}

func listAvailableProxyChannels(globalFlags *types.GlobalFlags, flags *listAvailableProxyChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := proxy.Proxy(&flags.ConnectionDetails, flags.Clientcert)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

