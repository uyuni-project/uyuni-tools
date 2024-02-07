
package admin

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "Provides methods to access and modify PAYG ssh connection data",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(setDetailsCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(listCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))

	return cmd
}
