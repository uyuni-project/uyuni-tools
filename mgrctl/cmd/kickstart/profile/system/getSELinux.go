package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getSELinuxFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func getSELinuxCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getSELinux",
		Short: "Retrieves the SELinux enforcing mode property of a kickstart
 profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getSELinuxFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getSELinux)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")

	return cmd
}

func getSELinux(globalFlags *types.GlobalFlags, flags *getSELinuxFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

