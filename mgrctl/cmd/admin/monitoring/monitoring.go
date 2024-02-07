
package monitoring

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitoring",
		Short: "Provides methods to manage the monitoring of the #product() server.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(enableCommand(globalFlags))
	cmd.AddCommand(disableCommand(globalFlags))
	cmd.AddCommand(getStatusCommand(globalFlags))

	return cmd
}
