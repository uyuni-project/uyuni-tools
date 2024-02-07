
package master

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "master",
		Short: "Contains methods to set up information about known-"masters", for use
 on the "slave" side of ISS",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(hasMasterCommand(globalFlags))
	cmd.AddCommand(makeDefaultCommand(globalFlags))
	cmd.AddCommand(setMasterOrgsCommand(globalFlags))
	cmd.AddCommand(addToMasterCommand(globalFlags))
	cmd.AddCommand(updateCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))
	cmd.AddCommand(getDefaultMasterCommand(globalFlags))
	cmd.AddCommand(getMasterByLabelCommand(globalFlags))
	cmd.AddCommand(getMasterOrgsCommand(globalFlags))
	cmd.AddCommand(setCaCertCommand(globalFlags))
	cmd.AddCommand(unsetDefaultMasterCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(mapToLocalCommand(globalFlags))
	cmd.AddCommand(getMasterCommand(globalFlags))
	cmd.AddCommand(getMastersCommand(globalFlags))

	return cmd
}
