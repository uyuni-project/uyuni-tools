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

type addActivationKeyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
	Key          string
}

func addActivationKeyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addActivationKey",
		Short: "Add an activation key association to the kickstart profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addActivationKeyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addActivationKey)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")
	cmd.Flags().String("Key", "", "the activation key")

	return cmd
}

func addActivationKey(globalFlags *types.GlobalFlags, flags *addActivationKeyFlags, cmd *cobra.Command, args []string) error {

res, err := keys.Keys(&flags.ConnectionDetails, flags.KsLabel, flags.Key)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

