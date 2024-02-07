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

type getSlaveByNameFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SlaveFqdn             string
}

func getSlaveByNameCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getSlaveByName",
		Short: "Find a Slave by specifying its Fully-Qualified Domain Name",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getSlaveByNameFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getSlaveByName)
		},
	}

	cmd.Flags().String("SlaveFqdn", "", "Domain-name of the desired Slave")

	return cmd
}

func getSlaveByName(globalFlags *types.GlobalFlags, flags *getSlaveByNameFlags, cmd *cobra.Command, args []string) error {

	res, err := slave.Slave(&flags.ConnectionDetails, flags.SlaveFqdn)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
