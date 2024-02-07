package maintenance

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/maintenance"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listScheduleNamesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listScheduleNamesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listScheduleNames",
		Short: "List Schedule Names visible to user",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listScheduleNamesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listScheduleNames)
		},
	}

	return cmd
}

func listScheduleNames(globalFlags *types.GlobalFlags, flags *listScheduleNamesFlags, cmd *cobra.Command, args []string) error {

	res, err := maintenance.Maintenance(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
