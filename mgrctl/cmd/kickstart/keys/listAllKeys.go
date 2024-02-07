package keys

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/keys"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listAllKeysFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listAllKeysCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAllKeys",
		Short: "list all keys for the org associated with the user logged into the
             given session",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAllKeysFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAllKeys)
		},
	}


	return cmd
}

func listAllKeys(globalFlags *types.GlobalFlags, flags *listAllKeysFlags, cmd *cobra.Command, args []string) error {

res, err := keys.Keys(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

