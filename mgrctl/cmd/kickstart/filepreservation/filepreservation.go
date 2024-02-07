
package filepreservation

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filepreservation",
		Short: "Provides methods to retrieve and manipulate kickstart file
 preservation lists.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(listAllFilePreservationsCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))

	return cmd
}
