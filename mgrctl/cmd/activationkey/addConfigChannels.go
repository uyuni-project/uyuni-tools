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

type addConfigChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Keys          []string
}

func addConfigChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addConfigChannels",
		Short: "Given a list of activation keys and configuration channels,
 this method adds given configuration channels to either the top or
 the bottom (whichever you specify) of an activation key's
 configuration channels list. The ordering of the configuration channels
 provided in the add list is maintained while adding.
 If one of the configuration channels in the 'add' list
 already exists in an activation key, the
 configuration  channel will be re-ranked to the appropriate place.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addConfigChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addConfigChannels)
		},
	}

	cmd.Flags().String("Keys", "", "$desc")

	return cmd
}

func addConfigChannels(globalFlags *types.GlobalFlags, flags *addConfigChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Keys)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

