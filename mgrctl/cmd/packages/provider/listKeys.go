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

type listKeysFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProviderName          string
}

func listKeysCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listKeys",
		Short: "List all security keys associated with a package provider.
 User executing the request must be a #product() administrator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listKeysFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listKeys)
		},
	}

	cmd.Flags().String("ProviderName", "", "The provider name")

	return cmd
}

func listKeys(globalFlags *types.GlobalFlags, flags *listKeysFlags, cmd *cobra.Command, args []string) error {

res, err := provider.Provider(&flags.ConnectionDetails, flags.ProviderName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

