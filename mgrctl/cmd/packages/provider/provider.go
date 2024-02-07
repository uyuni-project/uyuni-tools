
package provider

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider",
		Short: "Methods to retrieve information about Package Providers associated with
      packages.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(associateKeyCommand(globalFlags))
	cmd.AddCommand(listKeysCommand(globalFlags))
	cmd.AddCommand(listCommand(globalFlags))

	return cmd
}
