package audit

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/audit"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listSystemsByPatchStatusFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	CveIdentifier          string
}

func listSystemsByPatchStatusCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSystemsByPatchStatus",
		Short: "List visible systems with their patch status regarding a given CVE
 identifier. Please note that the query code relies on data that is pre-generated
 by the 'cve-server-channels' taskomatic job.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSystemsByPatchStatusFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSystemsByPatchStatus)
		},
	}

	cmd.Flags().String("CveIdentifier", "", "")

	return cmd
}

func listSystemsByPatchStatus(globalFlags *types.GlobalFlags, flags *listSystemsByPatchStatusFlags, cmd *cobra.Command, args []string) error {

res, err := audit.Audit(&flags.ConnectionDetails, flags.CveIdentifier)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

