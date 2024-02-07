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

type getKernelLivePatchFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func getKernelLivePatchCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getKernelLivePatch",
		Short: "Returns the currently active kernel live patching version relative to
 the running kernel version of the system, or empty string if live patching feature
 is not in use for the given system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getKernelLivePatchFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getKernelLivePatch)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getKernelLivePatch(globalFlags *types.GlobalFlags, flags *getKernelLivePatchFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

