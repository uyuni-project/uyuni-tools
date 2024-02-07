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

type setSELinuxFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func setSELinuxCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setSELinux",
		Short: "Sets the SELinux enforcing mode property of a kickstart profile
 so that a system created using this profile will be have
 the appropriate SELinux enforcing mode.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setSELinuxFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setSELinux)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")

	return cmd
}

func setSELinux(globalFlags *types.GlobalFlags, flags *setSELinuxFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

