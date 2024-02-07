
package recurring

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recurring",
		Short: "Provides methods to handle recurring actions for minions, system groups and organizations.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listByEntityCommand(globalFlags))
	cmd.AddCommand(lookupByIdCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))

	return cmd
}
