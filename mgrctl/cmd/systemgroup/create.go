package systemgroup

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/systemgroup"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type createFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name                  string
	Description           string
}

func createCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new system group.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, create)
		},
	}

	cmd.Flags().String("Name", "", "Name of the system group.")
	cmd.Flags().String("Description", "", "Description of the                  system group.")

	return cmd
}

func create(globalFlags *types.GlobalFlags, flags *createFlags, cmd *cobra.Command, args []string) error {

	res, err := systemgroup.Systemgroup(&flags.ConnectionDetails, flags.Name, flags.Description)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
