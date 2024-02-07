package search

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Methods to interface to package search capabilities in search server..",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(advancedWithChannelCommand(globalFlags))
	cmd.AddCommand(advancedWithActKeyCommand(globalFlags))
	cmd.AddCommand(advancedCommand(globalFlags))
	cmd.AddCommand(nameAndSummaryCommand(globalFlags))
	cmd.AddCommand(nameCommand(globalFlags))
	cmd.AddCommand(nameAndDescriptionCommand(globalFlags))

	return cmd
}
