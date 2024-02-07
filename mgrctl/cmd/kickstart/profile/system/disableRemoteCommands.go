package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type disableRemoteCommandsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func disableRemoteCommandsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disableRemoteCommands",
		Short: "Disables the remote command flag in a kickstart profile
 so that a system created using this profile
 will be capable of running remote commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags disableRemoteCommandsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, disableRemoteCommands)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")

	return cmd
}

func disableRemoteCommands(globalFlags *types.GlobalFlags, flags *disableRemoteCommandsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

