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

type enableRemoteCommandsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func enableRemoteCommandsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enableRemoteCommands",
		Short: "Enables the remote command flag in a kickstart profile
 so that a system created using this profile
 will be capable of running remote commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags enableRemoteCommandsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, enableRemoteCommands)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")

	return cmd
}

func enableRemoteCommands(globalFlags *types.GlobalFlags, flags *enableRemoteCommandsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

