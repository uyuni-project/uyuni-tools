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

type removeAnsiblePathFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	PathId                int
}

func removeAnsiblePathCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeAnsiblePath",
		Short: "Create ansible path",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeAnsiblePathFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeAnsiblePath)
		},
	}

	cmd.Flags().String("PathId", "", "path id")

	return cmd
}

func removeAnsiblePath(globalFlags *types.GlobalFlags, flags *removeAnsiblePathFlags, cmd *cobra.Command, args []string) error {

	res, err := ansible.Ansible(&flags.ConnectionDetails, flags.PathId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
