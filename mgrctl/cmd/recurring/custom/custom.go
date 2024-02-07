package custom

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom",
		Short: "Provides methods to handle recurring custom states for minions, system groups and organizations.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listAvailableCommand(globalFlags))
	cmd.AddCommand(updateCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))

	return cmd
}
