
package content

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "content",
		Short: "Provides the namespace for the content synchronization methods.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(addChannelsCommand(globalFlags))
	cmd.AddCommand(synchronizeSubscriptionsCommand(globalFlags))
	cmd.AddCommand(deleteCredentialsCommand(globalFlags))
	cmd.AddCommand(synchronizeRepositoriesCommand(globalFlags))
	cmd.AddCommand(listChannelsCommand(globalFlags))
	cmd.AddCommand(addCredentialsCommand(globalFlags))
	cmd.AddCommand(synchronizeProductsCommand(globalFlags))
	cmd.AddCommand(listProductsCommand(globalFlags))
	cmd.AddCommand(listCredentialsCommand(globalFlags))
	cmd.AddCommand(synchronizeChannelFamiliesCommand(globalFlags))
	cmd.AddCommand(addChannelCommand(globalFlags))
	cmd.AddCommand(synchronizeChannelsCommand(globalFlags))

	return cmd
}
