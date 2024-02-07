package packages

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/packages"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listProvidingChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Pid                   int
}

func listProvidingChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listProvidingChannels",
		Short: "List the channels that provide the a package.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listProvidingChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listProvidingChannels)
		},
	}

	cmd.Flags().String("Pid", "", "")

	return cmd
}

func listProvidingChannels(globalFlags *types.GlobalFlags, flags *listProvidingChannelsFlags, cmd *cobra.Command, args []string) error {

	res, err := packages.Packages(&flags.ConnectionDetails, flags.Pid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
