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

type listActivationKeysFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func listActivationKeysCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listActivationKeys",
		Short: "List the activation keys the system was registered with.  An empty
 list will be returned if an activation key was not used during registration.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listActivationKeysFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listActivationKeys)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listActivationKeys(globalFlags *types.GlobalFlags, flags *listActivationKeysFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

