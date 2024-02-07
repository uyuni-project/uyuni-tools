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

type removeServerGroupsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Key                   string
	ServerGroupIds        []int
}

func removeServerGroupsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeServerGroups",
		Short: "Remove server groups from an activation key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeServerGroupsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeServerGroups)
		},
	}

	cmd.Flags().String("Key", "", "")
	cmd.Flags().String("ServerGroupIds", "", "$desc")

	return cmd
}

func removeServerGroups(globalFlags *types.GlobalFlags, flags *removeServerGroupsFlags, cmd *cobra.Command, args []string) error {

	res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Key, flags.ServerGroupIds)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
