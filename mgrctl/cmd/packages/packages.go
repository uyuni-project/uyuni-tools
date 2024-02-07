
package packages

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "packages",
		Short: "Methods to retrieve information about the Packages contained
 within this server.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listChangelogCommand(globalFlags))
	cmd.AddCommand(getPackageCommand(globalFlags))
	cmd.AddCommand(findByNvreaCommand(globalFlags))
	cmd.AddCommand(removePackageCommand(globalFlags))
	cmd.AddCommand(listProvidingErrataCommand(globalFlags))
	cmd.AddCommand(listSourcePackagesCommand(globalFlags))
	cmd.AddCommand(listProvidingChannelsCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(removeSourcePackageCommand(globalFlags))
	cmd.AddCommand(listDependenciesCommand(globalFlags))
	cmd.AddCommand(listFilesCommand(globalFlags))
	cmd.AddCommand(getPackageUrlCommand(globalFlags))

	return cmd
}
