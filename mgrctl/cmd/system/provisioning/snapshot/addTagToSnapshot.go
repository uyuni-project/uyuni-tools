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

type addTagToSnapshotFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SnapId          int
	TagName          string
}

func addTagToSnapshotCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addTagToSnapshot",
		Short: "Adds tag to snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addTagToSnapshotFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addTagToSnapshot)
		},
	}

	cmd.Flags().String("SnapId", "", "ID of the snapshot")
	cmd.Flags().String("TagName", "", "Name of the snapshot tag")

	return cmd
}

func addTagToSnapshot(globalFlags *types.GlobalFlags, flags *addTagToSnapshotFlags, cmd *cobra.Command, args []string) error {

res, err := snapshot.Snapshot(&flags.ConnectionDetails, flags.SnapId, flags.TagName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

