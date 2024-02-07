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

type getSlavesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func getSlavesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getSlaves",
		Short: "Get all the Slaves this Master knows about",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getSlavesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getSlaves)
		},
	}


	return cmd
}

func getSlaves(globalFlags *types.GlobalFlags, flags *getSlavesFlags, cmd *cobra.Command, args []string) error {

res, err := slave.Slave(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

