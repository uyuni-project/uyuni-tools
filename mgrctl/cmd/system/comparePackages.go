package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type comparePackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid1                  int
	Sid2                  int
}

func comparePackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comparePackages",
		Short: "Compares the packages installed on two systems.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags comparePackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, comparePackages)
		},
	}

	cmd.Flags().String("Sid1", "", "")
	cmd.Flags().String("Sid2", "", "")

	return cmd
}

func comparePackages(globalFlags *types.GlobalFlags, flags *comparePackagesFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid1, flags.Sid2)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
