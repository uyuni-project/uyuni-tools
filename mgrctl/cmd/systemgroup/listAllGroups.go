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

type listAllGroupsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listAllGroupsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAllGroups",
		Short: "Retrieve a list of system groups that are accessible by the logged
      in user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAllGroupsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAllGroups)
		},
	}


	return cmd
}

func listAllGroups(globalFlags *types.GlobalFlags, flags *listAllGroupsFlags, cmd *cobra.Command, args []string) error {

res, err := systemgroup.Systemgroup(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

