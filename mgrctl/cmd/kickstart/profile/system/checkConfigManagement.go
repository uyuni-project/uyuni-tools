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

type checkConfigManagementFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func checkConfigManagementCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkConfigManagement",
		Short: "Check the configuration management status for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags checkConfigManagementFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, checkConfigManagement)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")

	return cmd
}

func checkConfigManagement(globalFlags *types.GlobalFlags, flags *checkConfigManagementFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

