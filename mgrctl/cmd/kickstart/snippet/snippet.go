package snippet

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snippet",
		Short: "Provides methods to create kickstart files",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(createOrUpdateCommand(globalFlags))
	cmd.AddCommand(listAllCommand(globalFlags))
	cmd.AddCommand(listCustomCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))
	cmd.AddCommand(listDefaultCommand(globalFlags))

	return cmd
}
