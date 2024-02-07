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

type listProxiesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listProxiesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listProxies",
		Short: "List the proxies within the user's organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listProxiesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listProxies)
		},
	}


	return cmd
}

func listProxies(globalFlags *types.GlobalFlags, flags *listProxiesFlags, cmd *cobra.Command, args []string) error {

res, err := proxy.Proxy(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

