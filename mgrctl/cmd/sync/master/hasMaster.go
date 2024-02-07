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

type hasMasterFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func hasMasterCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hasMaster",
		Short: "Check if this host is reading configuration from an ISS master.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags hasMasterFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, hasMaster)
		},
	}

	return cmd
}

func hasMaster(globalFlags *types.GlobalFlags, flags *hasMasterFlags, cmd *cobra.Command, args []string) error {

	res, err := master.Master(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
