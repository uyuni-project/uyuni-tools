package profile

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Provides methods to access and modify image profiles.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listImageProfileTypesCommand(globalFlags))
	cmd.AddCommand(getCustomValuesCommand(globalFlags))
	cmd.AddCommand(setDetailsCommand(globalFlags))
	cmd.AddCommand(listImageProfilesCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(deleteCustomValuesCommand(globalFlags))
	cmd.AddCommand(setCustomValuesCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))

	return cmd
}
