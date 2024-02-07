
package trusts

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trusts",
		Short: "Contains methods to access common organization trust information
 available from the web interface.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(removeTrustCommand(globalFlags))
	cmd.AddCommand(listTrustsCommand(globalFlags))
	cmd.AddCommand(listChannelsProvidedCommand(globalFlags))
	cmd.AddCommand(listSystemsAffectedCommand(globalFlags))
	cmd.AddCommand(listChannelsConsumedCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(addTrustCommand(globalFlags))
	cmd.AddCommand(listOrgsCommand(globalFlags))

	return cmd
}
