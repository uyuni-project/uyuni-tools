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

type listVirtualHostsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listVirtualHostsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listVirtualHosts",
		Short: "Lists the virtual hosts visible to the user",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listVirtualHostsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listVirtualHosts)
		},
	}

	return cmd
}

func listVirtualHosts(globalFlags *types.GlobalFlags, flags *listVirtualHostsFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
