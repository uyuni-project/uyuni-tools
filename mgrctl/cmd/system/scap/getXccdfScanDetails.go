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

type getXccdfScanDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Xid                   int
}

func getXccdfScanDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getXccdfScanDetails",
		Short: "Get details of given OpenSCAP XCCDF scan.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getXccdfScanDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getXccdfScanDetails)
		},
	}

	cmd.Flags().String("Xid", "", "ID of XCCDF scan.")

	return cmd
}

func getXccdfScanDetails(globalFlags *types.GlobalFlags, flags *getXccdfScanDetailsFlags, cmd *cobra.Command, args []string) error {

	res, err := scap.Scap(&flags.ConnectionDetails, flags.Xid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
