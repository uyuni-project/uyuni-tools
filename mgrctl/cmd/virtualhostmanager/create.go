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

type createFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	ModuleName          string
	Parameters          parameters
}

func createCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Virtual Host Manager from given arguments",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, create)
		},
	}

	cmd.Flags().String("Label", "", "Virtual Host Manager label")
	cmd.Flags().String("ModuleName", "", "the name of the Gatherer module")
	cmd.Flags().String("Parameters", "", "additional parameters (credentials, parameters for virtual-host-gatherer)")

	return cmd
}

func create(globalFlags *types.GlobalFlags, flags *createFlags, cmd *cobra.Command, args []string) error {

res, err := virtualhostmanager.Virtualhostmanager(&flags.ConnectionDetails, flags.Label, flags.ModuleName, flags.Parameters)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

