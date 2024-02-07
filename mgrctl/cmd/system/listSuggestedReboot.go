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

type listSuggestedRebootFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listSuggestedRebootCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSuggestedReboot",
		Short: "List systems that require reboot.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSuggestedRebootFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSuggestedReboot)
		},
	}

	return cmd
}

func listSuggestedReboot(globalFlags *types.GlobalFlags, flags *listSuggestedRebootFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
