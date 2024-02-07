package ansible

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/ansible"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listAnsiblePathsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ControlNodeId          int
}

func listAnsiblePathsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAnsiblePaths",
		Short: "List ansible paths for server (control node)",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAnsiblePathsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAnsiblePaths)
		},
	}

	cmd.Flags().String("ControlNodeId", "", "id of ansible control node server")

	return cmd
}

func listAnsiblePaths(globalFlags *types.GlobalFlags, flags *listAnsiblePathsFlags, cmd *cobra.Command, args []string) error {

res, err := ansible.Ansible(&flags.ConnectionDetails, flags.ControlNodeId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

