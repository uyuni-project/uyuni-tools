package filepreservation

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/filepreservation"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type deleteFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name                  string
}

func deleteCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a file preservation list.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, delete)
		},
	}

	cmd.Flags().String("Name", "", "name of the file list to delete")

	return cmd
}

func delete(globalFlags *types.GlobalFlags, flags *deleteFlags, cmd *cobra.Command, args []string) error {

	res, err := filepreservation.Filepreservation(&flags.ConnectionDetails, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
