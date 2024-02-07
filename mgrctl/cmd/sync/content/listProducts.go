package content

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/sync/content"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listProductsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listProductsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listProducts",
		Short: "List all accessible products.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listProductsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listProducts)
		},
	}


	return cmd
}

func listProducts(globalFlags *types.GlobalFlags, flags *listProductsFlags, cmd *cobra.Command, args []string) error {

res, err := content.Content(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

