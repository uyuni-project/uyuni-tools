package ansible

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ansible",
		Short: "Provides methods to manage Ansible systems",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(updateAnsiblePathCommand(globalFlags))
	cmd.AddCommand(createAnsiblePathCommand(globalFlags))
	cmd.AddCommand(fetchPlaybookContentsCommand(globalFlags))
	cmd.AddCommand(removeAnsiblePathCommand(globalFlags))
	cmd.AddCommand(introspectInventoryCommand(globalFlags))
	cmd.AddCommand(schedulePlaybookCommand(globalFlags))
	cmd.AddCommand(discoverPlaybooksCommand(globalFlags))
	cmd.AddCommand(listAnsiblePathsCommand(globalFlags))
	cmd.AddCommand(lookupAnsiblePathByIdCommand(globalFlags))

	return cmd
}
