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

type setLockStatusFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	LockStatus          bool
}

func setLockStatusCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setLockStatus",
		Short: "Set server lock status.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setLockStatusFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setLockStatus)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("LockStatus", "", "true to lock the system, false to unlock the system.")

	return cmd
}

func setLockStatus(globalFlags *types.GlobalFlags, flags *setLockStatusFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.LockStatus)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

