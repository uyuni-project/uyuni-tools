
package software

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "software",
		Short: "Provides methods to access and modify the software list
 associated with a kickstart profile.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(getSoftwareDetailsCommand(globalFlags))
	cmd.AddCommand(setSoftwareDetailsCommand(globalFlags))
	cmd.AddCommand(appendToSoftwareListCommand(globalFlags))
	cmd.AddCommand(getSoftwareListCommand(globalFlags))
	cmd.AddCommand(setSoftwareListCommand(globalFlags))

	return cmd
}
