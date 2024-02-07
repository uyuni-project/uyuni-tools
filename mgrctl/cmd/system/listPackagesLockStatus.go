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

type listPackagesLockStatusFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   string
}

func listPackagesLockStatusCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listPackagesLockStatus",
		Short: "List current package locks status.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listPackagesLockStatusFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listPackagesLockStatus)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listPackagesLockStatus(globalFlags *types.GlobalFlags, flags *listPackagesLockStatusFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
