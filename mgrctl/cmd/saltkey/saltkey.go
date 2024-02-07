
package saltkey

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "saltkey",
		Short: "Provides methods to manage salt keys",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(rejectCommand(globalFlags))
	cmd.AddCommand(deniedListCommand(globalFlags))
	cmd.AddCommand(rejectedListCommand(globalFlags))
	cmd.AddCommand(acceptedListCommand(globalFlags))
	cmd.AddCommand(pendingListCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))
	cmd.AddCommand(acceptCommand(globalFlags))

	return cmd
}
