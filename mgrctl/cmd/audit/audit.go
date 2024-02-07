
package audit

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Methods to audit systems.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listSystemsByPatchStatusCommand(globalFlags))
	cmd.AddCommand(listImagesByPatchStatusCommand(globalFlags))

	return cmd
}
