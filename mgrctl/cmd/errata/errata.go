package errata

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "errata",
		Short: "Provides methods to access and modify errata.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listKeywordsCommand(globalFlags))
	cmd.AddCommand(addPackagesCommand(globalFlags))
	cmd.AddCommand(applicableToChannelsCommand(globalFlags))
	cmd.AddCommand(bugzillaFixesCommand(globalFlags))
	cmd.AddCommand(removePackagesCommand(globalFlags))
	cmd.AddCommand(cloneAsyncCommand(globalFlags))
	cmd.AddCommand(listAffectedSystemsCommand(globalFlags))
	cmd.AddCommand(findByCveCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))
	cmd.AddCommand(listPackagesCommand(globalFlags))
	cmd.AddCommand(listCvesCommand(globalFlags))
	cmd.AddCommand(setDetailsCommand(globalFlags))
	cmd.AddCommand(publishCommand(globalFlags))
	cmd.AddCommand(cloneCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(cloneAsOriginalCommand(globalFlags))
	cmd.AddCommand(publishAsOriginalCommand(globalFlags))
	cmd.AddCommand(cloneAsOriginalAsyncCommand(globalFlags))

	return cmd
}
