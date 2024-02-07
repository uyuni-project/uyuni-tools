package activationkey

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/activationkey"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type disableConfigDeploymentFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Key          string
}

func disableConfigDeploymentCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disableConfigDeployment",
		Short: "Disable configuration file deployment for the specified activation key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags disableConfigDeploymentFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, disableConfigDeployment)
		},
	}

	cmd.Flags().String("Key", "", "")

	return cmd
}

func disableConfigDeployment(globalFlags *types.GlobalFlags, flags *disableConfigDeploymentFlags, cmd *cobra.Command, args []string) error {

res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Key)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

