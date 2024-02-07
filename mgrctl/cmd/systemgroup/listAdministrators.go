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

type listAdministratorsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemGroupName          string
}

func listAdministratorsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAdministrators",
		Short: "Returns the list of users who can administer the given group.
 Caller must be a system group admin or an organization administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAdministratorsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAdministrators)
		},
	}

	cmd.Flags().String("SystemGroupName", "", "")

	return cmd
}

func listAdministrators(globalFlags *types.GlobalFlags, flags *listAdministratorsFlags, cmd *cobra.Command, args []string) error {

res, err := systemgroup.Systemgroup(&flags.ConnectionDetails, flags.SystemGroupName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

