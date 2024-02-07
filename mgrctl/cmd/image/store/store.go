package store

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store",
		Short: "Provides methods to access and modify image stores.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listImageStoreTypesCommand(globalFlags))
	cmd.AddCommand(setDetailsCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(listImageStoresCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))

	return cmd
}
