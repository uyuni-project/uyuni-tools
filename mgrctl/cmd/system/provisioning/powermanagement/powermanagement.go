
package powermanagement

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "powermanagement",
		Short: "Provides methods to access and modify power management for systems.
 Some functions exist in 2 variants. Either with server id or with a name.
 The function with server id is useful when a system exists with a full profile.
 Everybody allowed to manage that system can execute these functions.
 The variant with name expects a cobbler system name prefix. These functions
 enhance the name by adding the org id of the user to limit access to systems
 from the own organization. Additionally Org Admin permissions are required to
 call these functions.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(rebootCommand(globalFlags))
	cmd.AddCommand(powerOnCommand(globalFlags))
	cmd.AddCommand(powerOffCommand(globalFlags))
	cmd.AddCommand(setDetailsCommand(globalFlags))
	cmd.AddCommand(listTypesCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(getStatusCommand(globalFlags))

	return cmd
}
