package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type createSystemProfileFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemName          string
	$param.getFlagName()          $param.getType()
}

func createSystemProfileCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createSystemProfile",
		Short: "Creates a system record in database for a system that is not registered.
 Either "hwAddress" or "hostname" prop must be specified in the "data" struct.
 If a system(s) matching given data exists, a SystemsExistFaultException is thrown which
 contains matching system IDs in its message.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createSystemProfileFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createSystemProfile)
		},
	}

	cmd.Flags().String("SystemName", "", "System name")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func createSystemProfile(globalFlags *types.GlobalFlags, flags *createSystemProfileFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.SystemName, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

