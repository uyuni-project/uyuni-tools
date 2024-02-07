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

type createAnsiblePathFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	$param.getFlagName()          $param.getType()
}

func createAnsiblePathCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createAnsiblePath",
		Short: "Create ansible path",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createAnsiblePathFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createAnsiblePath)
		},
	}

	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func createAnsiblePath(globalFlags *types.GlobalFlags, flags *createAnsiblePathFlags, cmd *cobra.Command, args []string) error {

res, err := ansible.Ansible(&flags.ConnectionDetails, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

