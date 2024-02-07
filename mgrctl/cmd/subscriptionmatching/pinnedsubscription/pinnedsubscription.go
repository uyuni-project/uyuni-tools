
package pinnedsubscription

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pinnedsubscription",
		Short: "Provides the namespace for operations on Pinned Subscriptions",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(listCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))

	return cmd
}
