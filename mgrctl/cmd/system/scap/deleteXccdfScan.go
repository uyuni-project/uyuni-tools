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

type deleteXccdfScanFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Xid          int
}

func deleteXccdfScanCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteXccdfScan",
		Short: "Delete OpenSCAP XCCDF Scan from the #product() database. Note that
 only those SCAP Scans can be deleted which have passed their retention period.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteXccdfScanFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteXccdfScan)
		},
	}

	cmd.Flags().String("Xid", "", "ID of XCCDF scan.")

	return cmd
}

func deleteXccdfScan(globalFlags *types.GlobalFlags, flags *deleteXccdfScanFlags, cmd *cobra.Command, args []string) error {

res, err := scap.Scap(&flags.ConnectionDetails, flags.Xid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

