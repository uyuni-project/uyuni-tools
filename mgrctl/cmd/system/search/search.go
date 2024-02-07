
package search

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Provides methods to perform system search requests using the search server.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(hostnameCommand(globalFlags))
	cmd.AddCommand(deviceDescriptionCommand(globalFlags))
	cmd.AddCommand(deviceVendorIdCommand(globalFlags))
	cmd.AddCommand(ipCommand(globalFlags))
	cmd.AddCommand(deviceDriverCommand(globalFlags))
	cmd.AddCommand(nameAndDescriptionCommand(globalFlags))
	cmd.AddCommand(uuidCommand(globalFlags))
	cmd.AddCommand(deviceIdCommand(globalFlags))

	return cmd
}
