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

type listAssignedSystemGroupsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label                 string
}

func listAssignedSystemGroupsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAssignedSystemGroups",
		Short: "Return a list of Groups where a given configuration channel is assigned to",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAssignedSystemGroupsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAssignedSystemGroups)
		},
	}

	cmd.Flags().String("Label", "", "label of the config channel to list assigned groups")

	return cmd
}

func listAssignedSystemGroups(globalFlags *types.GlobalFlags, flags *listAssignedSystemGroupsFlags, cmd *cobra.Command, args []string) error {

	res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
