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

type deleteSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sids                  []int
	CleanupType           string
}

func deleteSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteSystems",
		Short: "Delete systems given a list of system ids asynchronously.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteSystems)
		},
	}

	cmd.Flags().String("Sids", "", "$desc")
	cmd.Flags().String("CleanupType", "", "Possible values:  'FAIL_ON_CLEANUP_ERR' - fail in case of cleanup error,  'NO_CLEANUP' - do not cleanup, just delete,  'FORCE_DELETE' - Try cleanup first but delete server anyway in case of error")

	return cmd
}

func deleteSystems(globalFlags *types.GlobalFlags, flags *deleteSystemsFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sids, flags.CleanupType)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
