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

type getCpuFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func getCpuCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getCpu",
		Short: "Gets the CPU information of a system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getCpuFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getCpu)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getCpu(globalFlags *types.GlobalFlags, flags *getCpuFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

