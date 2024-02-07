
package access

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "access",
		Short: "Provides methods to retrieve and alter channel access restrictions.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(setOrgSharingCommand(globalFlags))
	cmd.AddCommand(disableUserRestrictionsCommand(globalFlags))
	cmd.AddCommand(getOrgSharingCommand(globalFlags))
	cmd.AddCommand(enableUserRestrictionsCommand(globalFlags))

	return cmd
}
