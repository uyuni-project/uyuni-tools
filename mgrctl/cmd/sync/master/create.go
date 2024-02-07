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

type createFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label                 string
}

func createCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Master, known to this Slave.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, create)
		},
	}

	cmd.Flags().String("Label", "", "Master's fully-qualified domain name")

	return cmd
}

func create(globalFlags *types.GlobalFlags, flags *createFlags, cmd *cobra.Command, args []string) error {

	res, err := master.Master(&flags.ConnectionDetails, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
