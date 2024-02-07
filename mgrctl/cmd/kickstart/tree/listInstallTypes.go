package tree

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/tree"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listInstallTypesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listInstallTypesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listInstallTypes",
		Short: "List the available kickstartable install types (rhel2,3,4,5 and
 fedora9+).",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listInstallTypesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listInstallTypes)
		},
	}


	return cmd
}

func listInstallTypes(globalFlags *types.GlobalFlags, flags *listInstallTypesFlags, cmd *cobra.Command, args []string) error {

res, err := tree.Tree(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

