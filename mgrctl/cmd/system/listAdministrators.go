package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listAdministratorsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func listAdministratorsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAdministrators",
		Short: "Returns a list of users which can administer the system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAdministratorsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAdministrators)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listAdministrators(globalFlags *types.GlobalFlags, flags *listAdministratorsFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
