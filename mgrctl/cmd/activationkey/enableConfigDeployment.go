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

type enableConfigDeploymentFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Key          string
}

func enableConfigDeploymentCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enableConfigDeployment",
		Short: "Enable configuration file deployment for the specified activation key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags enableConfigDeploymentFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, enableConfigDeployment)
		},
	}

	cmd.Flags().String("Key", "", "")

	return cmd
}

func enableConfigDeployment(globalFlags *types.GlobalFlags, flags *enableConfigDeploymentFlags, cmd *cobra.Command, args []string) error {

res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Key)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

