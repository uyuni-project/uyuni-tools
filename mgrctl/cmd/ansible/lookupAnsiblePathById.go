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

type lookupAnsiblePathByIdFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	PathId                int
}

func lookupAnsiblePathByIdCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lookupAnsiblePathById",
		Short: "Lookup ansible path by path id",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags lookupAnsiblePathByIdFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, lookupAnsiblePathById)
		},
	}

	cmd.Flags().String("PathId", "", "path id")

	return cmd
}

func lookupAnsiblePathById(globalFlags *types.GlobalFlags, flags *lookupAnsiblePathByIdFlags, cmd *cobra.Command, args []string) error {

	res, err := ansible.Ansible(&flags.ConnectionDetails, flags.PathId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
