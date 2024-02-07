
package snapshot

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Provides methods to access and delete system snapshots.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(addTagToSnapshotCommand(globalFlags))
	cmd.AddCommand(deleteSnapshotCommand(globalFlags))
	cmd.AddCommand(rollbackToTagCommand(globalFlags))
	cmd.AddCommand(listSnapshotConfigFilesCommand(globalFlags))
	cmd.AddCommand(listSnapshotsCommand(globalFlags))
	cmd.AddCommand(rollbackToSnapshotCommand(globalFlags))
	cmd.AddCommand(listSnapshotPackagesCommand(globalFlags))
	cmd.AddCommand(deleteSnapshotsCommand(globalFlags))

	return cmd
}
