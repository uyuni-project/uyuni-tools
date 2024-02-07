package scap

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scap",
		Short: "Provides methods to schedule SCAP scans and access the results.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(getXccdfScanRuleResultsCommand(globalFlags))
	cmd.AddCommand(deleteXccdfScanCommand(globalFlags))
	cmd.AddCommand(getXccdfScanDetailsCommand(globalFlags))
	cmd.AddCommand(scheduleXccdfScanCommand(globalFlags))
	cmd.AddCommand(listXccdfScansCommand(globalFlags))

	return cmd
}
