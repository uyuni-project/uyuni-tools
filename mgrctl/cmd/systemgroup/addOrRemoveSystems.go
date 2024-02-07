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

type addOrRemoveSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemGroupName       string
	ServerIds             []int
	Add                   bool
}

func addOrRemoveSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addOrRemoveSystems",
		Short: "Add/remove the given servers to a system group.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addOrRemoveSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addOrRemoveSystems)
		},
	}

	cmd.Flags().String("SystemGroupName", "", "")
	cmd.Flags().String("ServerIds", "", "$desc")
	cmd.Flags().String("Add", "", "True to add to the group,              False to remove.")

	return cmd
}

func addOrRemoveSystems(globalFlags *types.GlobalFlags, flags *addOrRemoveSystemsFlags, cmd *cobra.Command, args []string) error {

	res, err := systemgroup.Systemgroup(&flags.ConnectionDetails, flags.SystemGroupName, flags.ServerIds, flags.Add)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
