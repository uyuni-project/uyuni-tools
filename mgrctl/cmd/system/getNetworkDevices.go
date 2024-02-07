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

type getNetworkDevicesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func getNetworkDevicesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getNetworkDevices",
		Short: "Returns the network devices for the given server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getNetworkDevicesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getNetworkDevices)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getNetworkDevices(globalFlags *types.GlobalFlags, flags *getNetworkDevicesFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

