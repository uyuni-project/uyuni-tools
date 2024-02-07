
package delta

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delta",
		Short: "Provides methods to access and modify delta images.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(createDeltaImageCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(listDeltasCommand(globalFlags))

	return cmd
}
