package powermanagement

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/provisioning/powermanagement"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type powerOffFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	Name          string
}

func powerOffCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "powerOff",
		Short: "Execute power management action 'powerOff'",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags powerOffFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, powerOff)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Name", "", "")

	return cmd
}

func powerOff(globalFlags *types.GlobalFlags, flags *powerOffFlags, cmd *cobra.Command, args []string) error {

res, err := powermanagement.Powermanagement(&flags.ConnectionDetails, flags.Sid, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

