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

type listImagesByPatchStatusFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	CveIdentifier          string
}

func listImagesByPatchStatusCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listImagesByPatchStatus",
		Short: "List visible images with their patch status regarding a given CVE
 identifier. Please note that the query code relies on data that is pre-generated
 by the 'cve-server-channels' taskomatic job.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listImagesByPatchStatusFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listImagesByPatchStatus)
		},
	}

	cmd.Flags().String("CveIdentifier", "", "")

	return cmd
}

func listImagesByPatchStatus(globalFlags *types.GlobalFlags, flags *listImagesByPatchStatusFlags, cmd *cobra.Command, args []string) error {

res, err := audit.Audit(&flags.ConnectionDetails, flags.CveIdentifier)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

