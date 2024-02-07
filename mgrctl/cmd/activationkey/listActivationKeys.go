package activationkey

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/activationkey"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listActivationKeysFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listActivationKeysCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listActivationKeys",
		Short: "List activation keys that are visible to the
 user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listActivationKeysFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listActivationKeys)
		},
	}


	return cmd
}

func listActivationKeys(globalFlags *types.GlobalFlags, flags *listActivationKeysFlags, cmd *cobra.Command, args []string) error {

res, err := activationkey.Activationkey(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

