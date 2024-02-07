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

type disableConfigManagementFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func disableConfigManagementCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disableConfigManagement",
		Short: "Disables the configuration management flag in a kickstart profile
 so that a system created using this profile will be NOT be configuration capable.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags disableConfigManagementFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, disableConfigManagement)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")

	return cmd
}

func disableConfigManagement(globalFlags *types.GlobalFlags, flags *disableConfigManagementFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

