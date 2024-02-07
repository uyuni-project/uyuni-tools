package store

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/image/store"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listImageStoresFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listImageStoresCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listImageStores",
		Short: "List available image stores",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listImageStoresFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listImageStores)
		},
	}


	return cmd
}

func listImageStores(globalFlags *types.GlobalFlags, flags *listImageStoresFlags, cmd *cobra.Command, args []string) error {

res, err := store.Store(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

