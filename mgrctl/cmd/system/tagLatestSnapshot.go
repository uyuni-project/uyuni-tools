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

type tagLatestSnapshotFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	TagName          string
}

func tagLatestSnapshotCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tagLatestSnapshot",
		Short: "Tags latest system snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags tagLatestSnapshotFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, tagLatestSnapshot)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("TagName", "", "")

	return cmd
}

func tagLatestSnapshot(globalFlags *types.GlobalFlags, flags *tagLatestSnapshotFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.TagName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

