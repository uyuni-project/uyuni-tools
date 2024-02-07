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

type getMasterFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	MasterId          int
}

func getMasterCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getMaster",
		Short: "Find a Master by specifying its ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getMasterFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getMaster)
		},
	}

	cmd.Flags().String("MasterId", "", "ID of the desired Master")

	return cmd
}

func getMaster(globalFlags *types.GlobalFlags, flags *getMasterFlags, cmd *cobra.Command, args []string) error {

res, err := master.Master(&flags.ConnectionDetails, flags.MasterId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

