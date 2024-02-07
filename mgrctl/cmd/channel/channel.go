
package channel

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channel",
		Short: "Provides method to get back a list of Software Channels.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listMyChannelsCommand(globalFlags))
	cmd.AddCommand(listSoftwareChannelsCommand(globalFlags))
	cmd.AddCommand(listRetiredChannelsCommand(globalFlags))
	cmd.AddCommand(listSharedChannelsCommand(globalFlags))
	cmd.AddCommand(listManageableChannelsCommand(globalFlags))
	cmd.AddCommand(listVendorChannelsCommand(globalFlags))
	cmd.AddCommand(listAllChannelsCommand(globalFlags))
	cmd.AddCommand(listPopularChannelsCommand(globalFlags))

	return cmd
}
