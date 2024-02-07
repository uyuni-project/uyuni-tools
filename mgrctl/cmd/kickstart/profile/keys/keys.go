
package keys

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Provides methods to access and modify the list of activation keys
 associated with a kickstart profile.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(getActivationKeysCommand(globalFlags))
	cmd.AddCommand(addActivationKeyCommand(globalFlags))
	cmd.AddCommand(removeActivationKeyCommand(globalFlags))

	return cmd
}
