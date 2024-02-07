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

type setPrimaryInterfaceFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	InterfaceName          string
}

func setPrimaryInterfaceCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setPrimaryInterface",
		Short: "Sets new primary network interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setPrimaryInterfaceFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setPrimaryInterface)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("InterfaceName", "", "")

	return cmd
}

func setPrimaryInterface(globalFlags *types.GlobalFlags, flags *setPrimaryInterfaceFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.InterfaceName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

