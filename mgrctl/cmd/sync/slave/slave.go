
package slave

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slave",
		Short: "Contains methods to set up information about allowed-"slaves", for use
 on the "master" side of ISS",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(getAllowedOrgsCommand(globalFlags))
	cmd.AddCommand(setAllowedOrgsCommand(globalFlags))
	cmd.AddCommand(getSlaveCommand(globalFlags))
	cmd.AddCommand(getSlavesCommand(globalFlags))
	cmd.AddCommand(updateCommand(globalFlags))
	cmd.AddCommand(getSlaveByNameCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))

	return cmd
}
