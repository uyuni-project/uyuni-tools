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

type makeDefaultFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	MasterId          int
}

func makeDefaultCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "makeDefault",
		Short: "Make the specified Master the default for this Slave's inter-server-sync",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags makeDefaultFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, makeDefault)
		},
	}

	cmd.Flags().String("MasterId", "", "Id of the Master to make the default")

	return cmd
}

func makeDefault(globalFlags *types.GlobalFlags, flags *makeDefaultFlags, cmd *cobra.Command, args []string) error {

res, err := master.Master(&flags.ConnectionDetails, flags.MasterId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

