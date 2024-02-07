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

type deleteSystemFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ClientCert          string
	Sid          int
	CleanupType          string
}

func deleteSystemCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteSystem",
		Short: "Delete a system given its client certificate.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteSystemFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteSystem)
		},
	}

	cmd.Flags().String("ClientCert", "", "client certificate of the system")
	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("CleanupType", "", "Possible values:  'FAIL_ON_CLEANUP_ERR' - fail in case of cleanup error,  'NO_CLEANUP' - do not cleanup, just delete,  'FORCE_DELETE' - Try cleanup first but delete server anyway in case of error")

	return cmd
}

func deleteSystem(globalFlags *types.GlobalFlags, flags *deleteSystemFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.ClientCert, flags.Sid, flags.CleanupType)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

