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

type listGroupsWithNoAssociatedAdminsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listGroupsWithNoAssociatedAdminsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listGroupsWithNoAssociatedAdmins",
		Short: "Returns a list of system groups that do not have an administrator.
 (who is not an organization administrator, as they have implicit access to
 system groups) Caller must be an organization administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listGroupsWithNoAssociatedAdminsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listGroupsWithNoAssociatedAdmins)
		},
	}


	return cmd
}

func listGroupsWithNoAssociatedAdmins(globalFlags *types.GlobalFlags, flags *listGroupsWithNoAssociatedAdminsFlags, cmd *cobra.Command, args []string) error {

res, err := systemgroup.Systemgroup(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

