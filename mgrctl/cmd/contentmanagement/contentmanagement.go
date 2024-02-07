
package contentmanagement

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contentmanagement",
		Short: "Provides methods to access and modify Content Lifecycle Management related entities
 (Projects, Environments, Filters, Sources).",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(lookupProjectCommand(globalFlags))
	cmd.AddCommand(buildProjectCommand(globalFlags))
	cmd.AddCommand(updateFilterCommand(globalFlags))
	cmd.AddCommand(removeFilterCommand(globalFlags))
	cmd.AddCommand(createProjectCommand(globalFlags))
	cmd.AddCommand(lookupFilterCommand(globalFlags))
	cmd.AddCommand(lookupSourceCommand(globalFlags))
	cmd.AddCommand(listFiltersCommand(globalFlags))
	cmd.AddCommand(createAppStreamFiltersCommand(globalFlags))
	cmd.AddCommand(lookupEnvironmentCommand(globalFlags))
	cmd.AddCommand(removeEnvironmentCommand(globalFlags))
	cmd.AddCommand(listProjectEnvironmentsCommand(globalFlags))
	cmd.AddCommand(listProjectSourcesCommand(globalFlags))
	cmd.AddCommand(removeProjectCommand(globalFlags))
	cmd.AddCommand(detachSourceCommand(globalFlags))
	cmd.AddCommand(createFilterCommand(globalFlags))
	cmd.AddCommand(listProjectsCommand(globalFlags))
	cmd.AddCommand(listProjectFiltersCommand(globalFlags))
	cmd.AddCommand(attachFilterCommand(globalFlags))
	cmd.AddCommand(listFilterCriteriaCommand(globalFlags))
	cmd.AddCommand(updateProjectCommand(globalFlags))
	cmd.AddCommand(detachFilterCommand(globalFlags))
	cmd.AddCommand(createEnvironmentCommand(globalFlags))
	cmd.AddCommand(attachSourceCommand(globalFlags))
	cmd.AddCommand(updateEnvironmentCommand(globalFlags))
	cmd.AddCommand(promoteProjectCommand(globalFlags))

	return cmd
}
