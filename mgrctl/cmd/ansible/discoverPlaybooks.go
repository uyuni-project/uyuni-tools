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

type discoverPlaybooksFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	PathId                int
}

func discoverPlaybooksCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "discoverPlaybooks",
		Short: "Discover playbooks under given playbook path with given pathId",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags discoverPlaybooksFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, discoverPlaybooks)
		},
	}

	cmd.Flags().String("PathId", "", "path id")

	return cmd
}

func discoverPlaybooks(globalFlags *types.GlobalFlags, flags *discoverPlaybooksFlags, cmd *cobra.Command, args []string) error {

	res, err := ansible.Ansible(&flags.ConnectionDetails, flags.PathId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
