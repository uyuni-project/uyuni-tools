package tree

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tree",
		Short: "Provides methods to access and modify the kickstart trees.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listInstallTypesCommand(globalFlags))
	cmd.AddCommand(renameCommand(globalFlags))
	cmd.AddCommand(deleteTreeAndProfilesCommand(globalFlags))
	cmd.AddCommand(updateCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(listCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))

	return cmd
}
