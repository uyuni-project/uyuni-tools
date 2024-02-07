
package distchannel

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "distchannel",
		Short: "Provides methods to access and modify distribution channel information",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(setMapForOrgCommand(globalFlags))
	cmd.AddCommand(listDefaultMapsCommand(globalFlags))
	cmd.AddCommand(listMapsForOrgCommand(globalFlags))

	return cmd
}
