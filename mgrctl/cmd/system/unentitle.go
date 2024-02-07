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

type unentitleFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ClientCert            string
}

func unentitleCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unentitle",
		Short: "Unentitle the system completely",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags unentitleFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, unentitle)
		},
	}

	cmd.Flags().String("ClientCert", "", "client system id file")

	return cmd
}

func unentitle(globalFlags *types.GlobalFlags, flags *unentitleFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.ClientCert)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
