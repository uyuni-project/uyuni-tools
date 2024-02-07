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

type getDevicesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func getDevicesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getDevices",
		Short: "Gets a list of devices for a system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getDevicesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getDevices)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getDevices(globalFlags *types.GlobalFlags, flags *getDevicesFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
