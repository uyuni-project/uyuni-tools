
package monitoring

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitoring",
		Short: "Provides methods to access information about managed systems, applications and formulas which can be
 relevant for Prometheus monitoring",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listEndpointsCommand(globalFlags))

	return cmd
}
