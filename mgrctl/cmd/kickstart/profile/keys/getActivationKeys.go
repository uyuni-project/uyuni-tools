package keys

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile/keys"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getActivationKeysFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func getActivationKeysCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getActivationKeys",
		Short: "Lookup the activation keys associated with the kickstart
 profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getActivationKeysFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getActivationKeys)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")

	return cmd
}

func getActivationKeys(globalFlags *types.GlobalFlags, flags *getActivationKeysFlags, cmd *cobra.Command, args []string) error {

res, err := keys.Keys(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

