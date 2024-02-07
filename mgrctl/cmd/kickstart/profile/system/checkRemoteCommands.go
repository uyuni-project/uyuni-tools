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

type checkRemoteCommandsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func checkRemoteCommandsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkRemoteCommands",
		Short: "Check the remote commands status flag for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags checkRemoteCommandsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, checkRemoteCommands)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")

	return cmd
}

func checkRemoteCommands(globalFlags *types.GlobalFlags, flags *checkRemoteCommandsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

