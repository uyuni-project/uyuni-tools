
package schedule

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Methods to retrieve information about scheduled actions.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(failSystemActionCommand(globalFlags))
	cmd.AddCommand(listAllActionsCommand(globalFlags))
	cmd.AddCommand(archiveActionsCommand(globalFlags))
	cmd.AddCommand(listAllCompletedActionsCommand(globalFlags))
	cmd.AddCommand(listFailedActionsCommand(globalFlags))
	cmd.AddCommand(listAllArchivedActionsCommand(globalFlags))
	cmd.AddCommand(listCompletedActionsCommand(globalFlags))
	cmd.AddCommand(listFailedSystemsCommand(globalFlags))
	cmd.AddCommand(listCompletedSystemsCommand(globalFlags))
	cmd.AddCommand(listArchivedActionsCommand(globalFlags))
	cmd.AddCommand(deleteActionsCommand(globalFlags))
	cmd.AddCommand(listInProgressActionsCommand(globalFlags))
	cmd.AddCommand(listInProgressSystemsCommand(globalFlags))
	cmd.AddCommand(cancelActionsCommand(globalFlags))
	cmd.AddCommand(rescheduleActionsCommand(globalFlags))

	return cmd
}
