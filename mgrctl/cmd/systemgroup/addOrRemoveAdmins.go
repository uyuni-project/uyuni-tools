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

type addOrRemoveAdminsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemGroupName          string
	$param.getFlagName()          $param.getType()
	Add          int
}

func addOrRemoveAdminsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addOrRemoveAdmins",
		Short: "Add or remove administrators to/from the given group. #product() and
 Organization administrators are granted access to groups within their organization
 by default; therefore, users with those roles should not be included in the array
 provided. Caller must be an organization administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addOrRemoveAdminsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addOrRemoveAdmins)
		},
	}

	cmd.Flags().String("SystemGroupName", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("Add", "", "1 to add administrators, 0 to remove.")

	return cmd
}

func addOrRemoveAdmins(globalFlags *types.GlobalFlags, flags *addOrRemoveAdminsFlags, cmd *cobra.Command, args []string) error {

res, err := systemgroup.Systemgroup(&flags.ConnectionDetails, flags.SystemGroupName, flags.$param.getFlagName(), flags.Add)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

