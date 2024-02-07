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

type listProxyClientsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProxyId               int
}

func listProxyClientsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listProxyClients",
		Short: "List the clients directly connected to a given Proxy.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listProxyClientsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listProxyClients)
		},
	}

	cmd.Flags().String("ProxyId", "", "")

	return cmd
}

func listProxyClients(globalFlags *types.GlobalFlags, flags *listProxyClientsFlags, cmd *cobra.Command, args []string) error {

	res, err := proxy.Proxy(&flags.ConnectionDetails, flags.ProxyId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
