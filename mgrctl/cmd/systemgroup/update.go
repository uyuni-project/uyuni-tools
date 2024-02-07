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

type updateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemGroupName          string
	Description          string
}

func updateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing system group.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, update)
		},
	}

	cmd.Flags().String("SystemGroupName", "", "")
	cmd.Flags().String("Description", "", "")

	return cmd
}

func update(globalFlags *types.GlobalFlags, flags *updateFlags, cmd *cobra.Command, args []string) error {

res, err := systemgroup.Systemgroup(&flags.ConnectionDetails, flags.SystemGroupName, flags.Description)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

