package activationkey

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/activationkey"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	MinionId          string
	MachinePassword          string
	ActivationKey          string
}

func listChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listChannels",
		Short: "List the channels for the given activation key
 with temporary authentication tokens to access them.
 Authentication is done via a machine specific password.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listChannels)
		},
	}

	cmd.Flags().String("MinionId", "", "The id of the minion to authenticate with.")
	cmd.Flags().String("MachinePassword", "", "password specific to a machine.")
	cmd.Flags().String("ActivationKey", "", "activation key to use channels from.")

	return cmd
}

func listChannels(globalFlags *types.GlobalFlags, flags *listChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.MinionId, flags.MachinePassword, flags.ActivationKey)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

