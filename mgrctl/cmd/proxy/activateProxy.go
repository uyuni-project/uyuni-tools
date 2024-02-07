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

type activateProxyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Clientcert          string
	Version          string
}

func activateProxyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activateProxy",
		Short: "Activates the proxy identified by the given client
 certificate i.e. systemid file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags activateProxyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, activateProxy)
		},
	}

	cmd.Flags().String("Clientcert", "", "client certificate file")
	cmd.Flags().String("Version", "", "Version of proxy to be registered.")

	return cmd
}

func activateProxy(globalFlags *types.GlobalFlags, flags *activateProxyFlags, cmd *cobra.Command, args []string) error {

res, err := proxy.Proxy(&flags.ConnectionDetails, flags.Clientcert, flags.Version)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

