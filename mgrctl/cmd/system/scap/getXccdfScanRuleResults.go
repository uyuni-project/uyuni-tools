package scap

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/scap"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getXccdfScanRuleResultsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Xid          int
}

func getXccdfScanRuleResultsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getXccdfScanRuleResults",
		Short: "Return a full list of RuleResults for given OpenSCAP XCCDF scan.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getXccdfScanRuleResultsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getXccdfScanRuleResults)
		},
	}

	cmd.Flags().String("Xid", "", "ID of XCCDF scan.")

	return cmd
}

func getXccdfScanRuleResults(globalFlags *types.GlobalFlags, flags *getXccdfScanRuleResultsFlags, cmd *cobra.Command, args []string) error {

res, err := scap.Scap(&flags.ConnectionDetails, flags.Xid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

