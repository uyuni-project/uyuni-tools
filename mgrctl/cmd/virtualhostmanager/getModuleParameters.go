package virtualhostmanager

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/virtualhostmanager"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getModuleParametersFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ModuleName          string
}

func getModuleParametersCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getModuleParameters",
		Short: "Get a list of parameters for a virtual-host-gatherer module.
 It returns a map of parameters with their typical default values.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getModuleParametersFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getModuleParameters)
		},
	}

	cmd.Flags().String("ModuleName", "", "The name of the module")

	return cmd
}

func getModuleParameters(globalFlags *types.GlobalFlags, flags *getModuleParametersFlags, cmd *cobra.Command, args []string) error {

res, err := virtualhostmanager.Virtualhostmanager(&flags.ConnectionDetails, flags.ModuleName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

