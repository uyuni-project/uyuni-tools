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

type introspectInventoryFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	PathId          int
}

func introspectInventoryCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "introspectInventory",
		Short: "Introspect inventory under given inventory path with given pathId and return it in a structured way",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags introspectInventoryFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, introspectInventory)
		},
	}

	cmd.Flags().String("PathId", "", "path id")

	return cmd
}

func introspectInventory(globalFlags *types.GlobalFlags, flags *introspectInventoryFlags, cmd *cobra.Command, args []string) error {

res, err := ansible.Ansible(&flags.ConnectionDetails, flags.PathId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

