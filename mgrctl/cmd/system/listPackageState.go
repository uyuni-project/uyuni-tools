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

type listPackageStateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func listPackageStateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listPackageState",
		Short: "List possible migration targets for a system",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listPackageStateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listPackageState)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listPackageState(globalFlags *types.GlobalFlags, flags *listPackageStateFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

