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

type rollbackToTagFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	TagName          string
}

func rollbackToTagCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollbackToTag",
		Short: "Rollbacks server to snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags rollbackToTagFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, rollbackToTag)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("TagName", "", "Name of the snapshot tag")

	return cmd
}

func rollbackToTag(globalFlags *types.GlobalFlags, flags *rollbackToTagFlags, cmd *cobra.Command, args []string) error {

res, err := snapshot.Snapshot(&flags.ConnectionDetails, flags.Sid, flags.TagName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

