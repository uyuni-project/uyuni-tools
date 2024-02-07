package configchannel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/configchannel"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listSubscribedSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
}

func listSubscribedSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSubscribedSystems",
		Short: "Return a list of systems subscribed to a configuration channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSubscribedSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSubscribedSystems)
		},
	}

	cmd.Flags().String("Label", "", "label of the config channel to list subscribed systems")

	return cmd
}

func listSubscribedSystems(globalFlags *types.GlobalFlags, flags *listSubscribedSystemsFlags, cmd *cobra.Command, args []string) error {

res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

