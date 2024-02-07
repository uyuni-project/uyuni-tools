package virtualhostmanager

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/virtualhostmanager"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listVirtualHostManagersFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listVirtualHostManagersCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listVirtualHostManagers",
		Short: "Lists Virtual Host Managers visible to a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listVirtualHostManagersFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listVirtualHostManagers)
		},
	}


	return cmd
}

func listVirtualHostManagers(globalFlags *types.GlobalFlags, flags *listVirtualHostManagersFlags, cmd *cobra.Command, args []string) error {

res, err := virtualhostmanager.Virtualhostmanager(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

