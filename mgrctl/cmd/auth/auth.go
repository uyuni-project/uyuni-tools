
package auth

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "This namespace provides methods to authenticate with the system's
 management server.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(logoutCommand(globalFlags))
	cmd.AddCommand(loginCommand(globalFlags))

	return cmd
}
