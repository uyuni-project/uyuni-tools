package actionchain

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/actionchain"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type addConfigurationDeploymentFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChainLabel          string
	Sid          int
}

func addConfigurationDeploymentCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addConfigurationDeployment",
		Short: "Adds an action to deploy a configuration file to an Action Chain.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addConfigurationDeploymentFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addConfigurationDeployment)
		},
	}

	cmd.Flags().String("ChainLabel", "", "Label of the chain")
	cmd.Flags().String("Sid", "", "System ID")

	return cmd
}

func addConfigurationDeployment(globalFlags *types.GlobalFlags, flags *addConfigurationDeploymentFlags, cmd *cobra.Command, args []string) error {

res, err := actionchain.Actionchain(&flags.ConnectionDetails, flags.ChainLabel, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

