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

type addPackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Key          string
	$param.getFlagName()          $param.getType()
}

func addPackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addPackages",
		Short: "Add packages to an activation key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addPackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addPackages)
		},
	}

	cmd.Flags().String("Key", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func addPackages(globalFlags *types.GlobalFlags, flags *addPackagesFlags, cmd *cobra.Command, args []string) error {

res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Key, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

