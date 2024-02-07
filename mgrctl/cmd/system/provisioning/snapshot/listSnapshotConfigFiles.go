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

type listSnapshotConfigFilesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SnapId          int
}

func listSnapshotConfigFilesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSnapshotConfigFiles",
		Short: "List the config files associated with a snapshot.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSnapshotConfigFilesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSnapshotConfigFiles)
		},
	}

	cmd.Flags().String("SnapId", "", "")

	return cmd
}

func listSnapshotConfigFiles(globalFlags *types.GlobalFlags, flags *listSnapshotConfigFilesFlags, cmd *cobra.Command, args []string) error {

res, err := snapshot.Snapshot(&flags.ConnectionDetails, flags.SnapId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

