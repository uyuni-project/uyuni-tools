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

type listSnapshotsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	StartDate          $type
	EndDate          $type
}

func listSnapshotsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSnapshots",
		Short: "List snapshots for a given system.
 A user may optionally provide a start and end date to narrow the snapshots that
 will be listed.  For example,
 
 If the user provides startDate only, all snapshots created either on or after
 the date provided will be returned.
 If user provides startDate and endDate, all snapshots created on or between the
 dates provided will be returned.
 If the user doesn't provide a startDate and endDate, all snapshots associated
 with the server will be returned.
 ",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSnapshotsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSnapshots)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("StartDate", "", "")
	cmd.Flags().String("EndDate", "", "")

	return cmd
}

func listSnapshots(globalFlags *types.GlobalFlags, flags *listSnapshotsFlags, cmd *cobra.Command, args []string) error {

res, err := snapshot.Snapshot(&flags.ConnectionDetails, flags.Sid, flags.StartDate, flags.EndDate)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

