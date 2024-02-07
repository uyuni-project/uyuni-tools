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

type obtainReactivationKeyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	ClientCert          string
}

func obtainReactivationKeyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "obtainReactivationKey",
		Short: "Obtains a reactivation key for this server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags obtainReactivationKeyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, obtainReactivationKey)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("ClientCert", "", "client certificate of the system")

	return cmd
}

func obtainReactivationKey(globalFlags *types.GlobalFlags, flags *obtainReactivationKeyFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.ClientCert)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

