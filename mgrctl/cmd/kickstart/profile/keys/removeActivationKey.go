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

type removeActivationKeyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
	Key          string
}

func removeActivationKeyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeActivationKey",
		Short: "Remove an activation key association from the kickstart profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeActivationKeyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeActivationKey)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")
	cmd.Flags().String("Key", "", "the activation key")

	return cmd
}

func removeActivationKey(globalFlags *types.GlobalFlags, flags *removeActivationKeyFlags, cmd *cobra.Command, args []string) error {

res, err := keys.Keys(&flags.ConnectionDetails, flags.KsLabel, flags.Key)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

