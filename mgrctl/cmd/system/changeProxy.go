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

type changeProxyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sids          []int
	ProxyId          int
}

func changeProxyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changeProxy",
		Short: "Connect given systems to another proxy.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags changeProxyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, changeProxy)
		},
	}

	cmd.Flags().String("Sids", "", "$desc")
	cmd.Flags().String("ProxyId", "", "")

	return cmd
}

func changeProxy(globalFlags *types.GlobalFlags, flags *changeProxyFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sids, flags.ProxyId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

