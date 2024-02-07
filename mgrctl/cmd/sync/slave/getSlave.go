package slave

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/sync/slave"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getSlaveFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SlaveId          int
}

func getSlaveCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getSlave",
		Short: "Find a Slave by specifying its ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getSlaveFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getSlave)
		},
	}

	cmd.Flags().String("SlaveId", "", "ID of the desired Slave")

	return cmd
}

func getSlave(globalFlags *types.GlobalFlags, flags *getSlaveFlags, cmd *cobra.Command, args []string) error {

res, err := slave.Slave(&flags.ConnectionDetails, flags.SlaveId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

