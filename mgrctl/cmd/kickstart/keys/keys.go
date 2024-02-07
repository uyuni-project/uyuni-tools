package keys

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Provides methods to manipulate kickstart keys.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listAllKeysCommand(globalFlags))
	cmd.AddCommand(updateCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))

	return cmd
}
