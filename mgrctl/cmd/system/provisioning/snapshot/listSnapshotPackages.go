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

type listSnapshotPackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SnapId          int
}

func listSnapshotPackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSnapshotPackages",
		Short: "List the packages associated with a snapshot.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSnapshotPackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSnapshotPackages)
		},
	}

	cmd.Flags().String("SnapId", "", "")

	return cmd
}

func listSnapshotPackages(globalFlags *types.GlobalFlags, flags *listSnapshotPackagesFlags, cmd *cobra.Command, args []string) error {

res, err := snapshot.Snapshot(&flags.ConnectionDetails, flags.SnapId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

