package master

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/sync/master"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type updateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	MasterId          int
	Label          string
}

func updateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates the label of the specified Master",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, update)
		},
	}

	cmd.Flags().String("MasterId", "", "ID of the Master to update")
	cmd.Flags().String("Label", "", "Desired new label")

	return cmd
}

func update(globalFlags *types.GlobalFlags, flags *updateFlags, cmd *cobra.Command, args []string) error {

res, err := master.Master(&flags.ConnectionDetails, flags.MasterId, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

