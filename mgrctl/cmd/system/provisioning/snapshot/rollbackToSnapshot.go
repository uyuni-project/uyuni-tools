package snapshot

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/provisioning/snapshot"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type rollbackToSnapshotFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	SnapId          int
}

func rollbackToSnapshotCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollbackToSnapshot",
		Short: "Rollbacks server to snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags rollbackToSnapshotFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, rollbackToSnapshot)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("SnapId", "", "ID of the snapshot")

	return cmd
}

func rollbackToSnapshot(globalFlags *types.GlobalFlags, flags *rollbackToSnapshotFlags, cmd *cobra.Command, args []string) error {

res, err := snapshot.Snapshot(&flags.ConnectionDetails, flags.Sid, flags.SnapId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

