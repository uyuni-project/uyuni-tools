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

type listImageStoreTypesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listImageStoreTypesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listImageStoreTypes",
		Short: "List available image store types",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listImageStoreTypesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listImageStoreTypes)
		},
	}


	return cmd
}

func listImageStoreTypes(globalFlags *types.GlobalFlags, flags *listImageStoreTypesFlags, cmd *cobra.Command, args []string) error {

res, err := store.Store(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

