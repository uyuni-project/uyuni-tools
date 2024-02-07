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

type getMastersFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func getMastersCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getMasters",
		Short: "Get all the Masters this Slave knows about",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getMastersFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getMasters)
		},
	}


	return cmd
}

func getMasters(globalFlags *types.GlobalFlags, flags *getMastersFlags, cmd *cobra.Command, args []string) error {

res, err := master.Master(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

