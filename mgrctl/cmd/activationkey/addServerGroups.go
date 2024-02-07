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

type addServerGroupsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Key                   string
	ServerGroupIds        []int
}

func addServerGroupsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addServerGroups",
		Short: "Add server groups to an activation key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addServerGroupsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addServerGroups)
		},
	}

	cmd.Flags().String("Key", "", "")
	cmd.Flags().String("ServerGroupIds", "", "$desc")

	return cmd
}

func addServerGroups(globalFlags *types.GlobalFlags, flags *addServerGroupsFlags, cmd *cobra.Command, args []string) error {

	res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Key, flags.ServerGroupIds)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
