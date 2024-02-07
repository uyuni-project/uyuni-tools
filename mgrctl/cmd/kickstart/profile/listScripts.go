package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listScriptsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func listScriptsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listScripts",
		Short: "List the pre and post scripts for a kickstart profile
 in the order they will run during the kickstart.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listScriptsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listScripts)
		},
	}

	cmd.Flags().String("KsLabel", "", "The label of the kickstart")

	return cmd
}

func listScripts(globalFlags *types.GlobalFlags, flags *listScriptsFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

