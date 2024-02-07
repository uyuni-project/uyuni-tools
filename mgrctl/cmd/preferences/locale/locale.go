package locale

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locale",
		Short: "Provides methods to access and modify user locale information",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(setTimeZoneCommand(globalFlags))
	cmd.AddCommand(listTimeZonesCommand(globalFlags))
	cmd.AddCommand(setLocaleCommand(globalFlags))
	cmd.AddCommand(listLocalesCommand(globalFlags))

	return cmd
}
