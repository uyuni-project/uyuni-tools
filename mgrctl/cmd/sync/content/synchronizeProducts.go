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

type synchronizeProductsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func synchronizeProductsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "synchronizeProducts",
		Short: "Synchronize SUSE products between the Customer Center
             and the #product() database.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags synchronizeProductsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, synchronizeProducts)
		},
	}


	return cmd
}

func synchronizeProducts(globalFlags *types.GlobalFlags, flags *synchronizeProductsFlags, cmd *cobra.Command, args []string) error {

res, err := content.Content(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

