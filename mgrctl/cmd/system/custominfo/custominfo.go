
package custominfo

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custominfo",
		Short: "Provides methods to access and modify custom system information.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listAllKeysCommand(globalFlags))
	cmd.AddCommand(updateKeyCommand(globalFlags))
	cmd.AddCommand(createKeyCommand(globalFlags))
	cmd.AddCommand(deleteKeyCommand(globalFlags))

	return cmd
}
