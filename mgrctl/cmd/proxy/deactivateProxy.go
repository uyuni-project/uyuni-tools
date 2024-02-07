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

type deactivateProxyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Clientcert          string
}

func deactivateProxyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deactivateProxy",
		Short: "Deactivates the proxy identified by the given client
 certificate i.e. systemid file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deactivateProxyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deactivateProxy)
		},
	}

	cmd.Flags().String("Clientcert", "", "client certificate file")

	return cmd
}

func deactivateProxy(globalFlags *types.GlobalFlags, flags *deactivateProxyFlags, cmd *cobra.Command, args []string) error {

res, err := proxy.Proxy(&flags.ConnectionDetails, flags.Clientcert)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

