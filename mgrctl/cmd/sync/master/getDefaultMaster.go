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

type getDefaultMasterFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func getDefaultMasterCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getDefaultMaster",
		Short: "Return the current default-Master for this Slave",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getDefaultMasterFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getDefaultMaster)
		},
	}

	return cmd
}

func getDefaultMaster(globalFlags *types.GlobalFlags, flags *getDefaultMasterFlags, cmd *cobra.Command, args []string) error {

	res, err := master.Master(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
