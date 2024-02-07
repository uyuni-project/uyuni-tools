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

type getRunningKernelFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func getRunningKernelCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getRunningKernel",
		Short: "Returns the running kernel of the given system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getRunningKernelFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getRunningKernel)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getRunningKernel(globalFlags *types.GlobalFlags, flags *getRunningKernelFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
