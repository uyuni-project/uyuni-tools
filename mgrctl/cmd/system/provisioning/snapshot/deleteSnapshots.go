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

type deleteSnapshotsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	StartDate          $type
	EndDate          $type
	Sid          int
	Sid          int
}

func deleteSnapshotsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteSnapshots",
		Short: "Deletes all snapshots across multiple systems based on the given date
 criteria.  For example,
 
 If the user provides startDate only, all snapshots created either on or after
 the date provided will be removed.
 If user provides startDate and endDate, all snapshots created on or between the
 dates provided will be removed.
 If the user doesn't provide a startDate and endDate, all snapshots will be
 removed.
 ",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteSnapshotsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteSnapshots)
		},
	}

	cmd.Flags().String("StartDate", "", "")
	cmd.Flags().String("EndDate", "", "")
	cmd.Flags().String("Sid", "", "ID of system to delete snapshots for")
	cmd.Flags().String("Sid", "", "ID of system to delete          snapshots for")

	return cmd
}

func deleteSnapshots(globalFlags *types.GlobalFlags, flags *deleteSnapshotsFlags, cmd *cobra.Command, args []string) error {

res, err := snapshot.Snapshot(&flags.ConnectionDetails, flags.StartDate, flags.EndDate, flags.Sid, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

