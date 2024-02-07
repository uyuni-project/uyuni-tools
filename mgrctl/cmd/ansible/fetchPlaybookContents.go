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

type fetchPlaybookContentsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	PathId                int
	PlaybookRelPath       string
}

func fetchPlaybookContentsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fetchPlaybookContents",
		Short: "Fetch the playbook content from the control node using a synchronous salt call.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags fetchPlaybookContentsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, fetchPlaybookContents)
		},
	}

	cmd.Flags().String("PathId", "", "playbook path id")
	cmd.Flags().String("PlaybookRelPath", "", "relative path of playbook (inside path specified by pathId)")

	return cmd
}

func fetchPlaybookContents(globalFlags *types.GlobalFlags, flags *fetchPlaybookContentsFlags, cmd *cobra.Command, args []string) error {

	res, err := ansible.Ansible(&flags.ConnectionDetails, flags.PathId, flags.PlaybookRelPath)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
