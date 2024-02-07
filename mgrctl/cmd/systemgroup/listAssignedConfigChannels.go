package systemgroup

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/systemgroup"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listAssignedConfigChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemGroupName       string
}

func listAssignedConfigChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAssignedConfigChannels",
		Short: "List all Configuration Channels assigned to a system group",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAssignedConfigChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAssignedConfigChannels)
		},
	}

	cmd.Flags().String("SystemGroupName", "", "")

	return cmd
}

func listAssignedConfigChannels(globalFlags *types.GlobalFlags, flags *listAssignedConfigChannelsFlags, cmd *cobra.Command, args []string) error {

	res, err := systemgroup.Systemgroup(&flags.ConnectionDetails, flags.SystemGroupName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
