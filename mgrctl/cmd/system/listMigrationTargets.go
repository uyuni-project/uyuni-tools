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

type listMigrationTargetsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	ExcludeTargetWhereMissingSuccessors          bool
}

func listMigrationTargetsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listMigrationTargets",
		Short: "List possible migration targets for a system",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listMigrationTargetsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listMigrationTargets)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("ExcludeTargetWhereMissingSuccessors", "", "")

	return cmd
}

func listMigrationTargets(globalFlags *types.GlobalFlags, flags *listMigrationTargetsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.ExcludeTargetWhereMissingSuccessors)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

