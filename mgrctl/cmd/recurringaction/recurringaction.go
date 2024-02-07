
package recurringaction

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recurringaction",
		Short: "Provides methods to handle recurring actions for minions, system groups and organizations.
 
 Deprecated - This namespace will be removed in a future API version. To work with recurring actions,
 please check out the newer 'recurring' namespace.",
	}

	api.AddAPIFlags(cmd, false)


	return cmd
}
