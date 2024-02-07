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

type listProvidingErrataFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Pid                   int
}

func listProvidingErrataCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listProvidingErrata",
		Short: "List the errata providing the a package.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listProvidingErrataFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listProvidingErrata)
		},
	}

	cmd.Flags().String("Pid", "", "")

	return cmd
}

func listProvidingErrata(globalFlags *types.GlobalFlags, flags *listProvidingErrataFlags, cmd *cobra.Command, args []string) error {

	res, err := packages.Packages(&flags.ConnectionDetails, flags.Pid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
