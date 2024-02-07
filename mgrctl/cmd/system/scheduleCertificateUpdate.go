package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type scheduleCertificateUpdateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	EarliestOccurrence          $date
}

func scheduleCertificateUpdateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleCertificateUpdate",
		Short: "Schedule update of client certificate",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleCertificateUpdateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleCertificateUpdate)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("EarliestOccurrence", "", "")

	return cmd
}

func scheduleCertificateUpdate(globalFlags *types.GlobalFlags, flags *scheduleCertificateUpdateFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.EarliestOccurrence)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

