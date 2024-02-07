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

type unsetDefaultMasterFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func unsetDefaultMasterCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unsetDefaultMaster",
		Short: "Make this slave have no default Master for inter-server-sync",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags unsetDefaultMasterFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, unsetDefaultMaster)
		},
	}

	return cmd
}

func unsetDefaultMaster(globalFlags *types.GlobalFlags, flags *unsetDefaultMasterFlags, cmd *cobra.Command, args []string) error {

	res, err := master.Master(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
