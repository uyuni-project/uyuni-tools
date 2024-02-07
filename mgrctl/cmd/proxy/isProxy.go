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

type isProxyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Clientcert          string
}

func isProxyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "isProxy",
		Short: "Test, if the system identified by the given client
 certificate i.e. systemid file, is proxy.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags isProxyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, isProxy)
		},
	}

	cmd.Flags().String("Clientcert", "", "client certificate file")

	return cmd
}

func isProxy(globalFlags *types.GlobalFlags, flags *isProxyFlags, cmd *cobra.Command, args []string) error {

res, err := proxy.Proxy(&flags.ConnectionDetails, flags.Clientcert)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

