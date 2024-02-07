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

type checkConfigDeploymentFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Key          string
}

func checkConfigDeploymentCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkConfigDeployment",
		Short: "Check configuration file deployment status for the
 activation key specified.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags checkConfigDeploymentFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, checkConfigDeployment)
		},
	}

	cmd.Flags().String("Key", "", "")

	return cmd
}

func checkConfigDeployment(globalFlags *types.GlobalFlags, flags *checkConfigDeploymentFlags, cmd *cobra.Command, args []string) error {

res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Key)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

