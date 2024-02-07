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

type listConfigChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Key          string
}

func listConfigChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listConfigChannels",
		Short: "List configuration channels
 associated to an activation key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listConfigChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listConfigChannels)
		},
	}

	cmd.Flags().String("Key", "", "")

	return cmd
}

func listConfigChannels(globalFlags *types.GlobalFlags, flags *listConfigChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Key)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

