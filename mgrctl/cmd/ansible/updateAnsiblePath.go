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

type updateAnsiblePathFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	PathId          int
	$param.getFlagName()          $param.getType()
}

func updateAnsiblePathCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateAnsiblePath",
		Short: "Create ansible path",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateAnsiblePathFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, updateAnsiblePath)
		},
	}

	cmd.Flags().String("PathId", "", "path id")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func updateAnsiblePath(globalFlags *types.GlobalFlags, flags *updateAnsiblePathFlags, cmd *cobra.Command, args []string) error {

res, err := ansible.Ansible(&flags.ConnectionDetails, flags.PathId, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

