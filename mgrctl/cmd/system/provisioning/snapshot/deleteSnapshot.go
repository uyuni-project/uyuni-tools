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

type deleteSnapshotFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SnapId                int
}

func deleteSnapshotCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteSnapshot",
		Short: "Deletes a snapshot with the given snapshot id",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteSnapshotFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteSnapshot)
		},
	}

	cmd.Flags().String("SnapId", "", "ID of snapshot to delete")

	return cmd
}

func deleteSnapshot(globalFlags *types.GlobalFlags, flags *deleteSnapshotFlags, cmd *cobra.Command, args []string) error {

	res, err := snapshot.Snapshot(&flags.ConnectionDetails, flags.SnapId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
