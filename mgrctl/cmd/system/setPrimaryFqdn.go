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

type setPrimaryFqdnFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
	Fqdn                  string
}

func setPrimaryFqdnCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setPrimaryFqdn",
		Short: "Sets new primary FQDN",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setPrimaryFqdnFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setPrimaryFqdn)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Fqdn", "", "")

	return cmd
}

func setPrimaryFqdn(globalFlags *types.GlobalFlags, flags *setPrimaryFqdnFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.Fqdn)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
