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

type removePackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Key          string
	$param.getFlagName()          $param.getType()
}

func removePackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removePackages",
		Short: "Remove package names from an activation key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removePackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removePackages)
		},
	}

	cmd.Flags().String("Key", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func removePackages(globalFlags *types.GlobalFlags, flags *removePackagesFlags, cmd *cobra.Command, args []string) error {

res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Key, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

