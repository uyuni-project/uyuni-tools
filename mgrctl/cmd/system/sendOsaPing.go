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

type sendOsaPingFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ServerId          int
}

func sendOsaPingCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sendOsaPing",
		Short: "send a ping to a system using OSA",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags sendOsaPingFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, sendOsaPing)
		},
	}

	cmd.Flags().String("ServerId", "", "")

	return cmd
}

func sendOsaPing(globalFlags *types.GlobalFlags, flags *sendOsaPingFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.ServerId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

