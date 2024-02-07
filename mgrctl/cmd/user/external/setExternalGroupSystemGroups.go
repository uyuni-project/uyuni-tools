package external

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/user/external"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setExternalGroupSystemGroupsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
}

func setExternalGroupSystemGroupsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setExternalGroupSystemGroups",
		Short: "Update the server groups for an external group. Replace previously set
 server groups with the ones passed in here. Can only be called by an org_admin.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setExternalGroupSystemGroupsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setExternalGroupSystemGroups)
		},
	}

	cmd.Flags().String("Name", "", "Name of the external group.")

	return cmd
}

func setExternalGroupSystemGroups(globalFlags *types.GlobalFlags, flags *setExternalGroupSystemGroupsFlags, cmd *cobra.Command, args []string) error {

res, err := external.External(&flags.ConnectionDetails, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

