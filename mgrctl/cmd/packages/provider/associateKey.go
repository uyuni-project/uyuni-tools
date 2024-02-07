package provider

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/packages/provider"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type associateKeyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProviderName          string
	Key          string
	Type          string
}

func associateKeyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "associateKey",
		Short: "Associate a package security key and with the package provider.
      If the provider or key doesn't exist, it is created. User executing the
      request must be a #product() administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags associateKeyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, associateKey)
		},
	}

	cmd.Flags().String("ProviderName", "", "The provider name")
	cmd.Flags().String("Key", "", "The actual key")
	cmd.Flags().String("Type", "", "The type of the key. Currently, only 'gpg' is supported")

	return cmd
}

func associateKey(globalFlags *types.GlobalFlags, flags *associateKeyFlags, cmd *cobra.Command, args []string) error {

res, err := provider.Provider(&flags.ConnectionDetails, flags.ProviderName, flags.Key, flags.Type)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

