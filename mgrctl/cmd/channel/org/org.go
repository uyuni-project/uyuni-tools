
package org

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org",
		Short: "Provides methods to retrieve and alter organization trust
 relationships for a channel.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(disableAccessCommand(globalFlags))
	cmd.AddCommand(enableAccessCommand(globalFlags))
	cmd.AddCommand(listCommand(globalFlags))

	return cmd
}
