package master

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/sync/master"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setCaCertFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	MasterId          int
	CaCertFilename          string
}

func setCaCertCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setCaCert",
		Short: "Set the CA-CERT filename for specified Master on this Slave",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setCaCertFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setCaCert)
		},
	}

	cmd.Flags().String("MasterId", "", "ID of the Master to affect")
	cmd.Flags().String("CaCertFilename", "", "path to specified Master's CA cert")

	return cmd
}

func setCaCert(globalFlags *types.GlobalFlags, flags *setCaCertFlags, cmd *cobra.Command, args []string) error {

res, err := master.Master(&flags.ConnectionDetails, flags.MasterId, flags.CaCertFilename)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

