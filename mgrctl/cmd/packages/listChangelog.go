package packages

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/packages"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listChangelogFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Pid                   int
}

func listChangelogCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listChangelog",
		Short: "List the change log for a package.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listChangelogFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listChangelog)
		},
	}

	cmd.Flags().String("Pid", "", "")

	return cmd
}

func listChangelog(globalFlags *types.GlobalFlags, flags *listChangelogFlags, cmd *cobra.Command, args []string) error {

	res, err := packages.Packages(&flags.ConnectionDetails, flags.Pid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
