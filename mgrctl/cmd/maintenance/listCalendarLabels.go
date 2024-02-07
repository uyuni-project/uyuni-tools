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

type listCalendarLabelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listCalendarLabelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listCalendarLabels",
		Short: "List schedule names visible to user",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listCalendarLabelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listCalendarLabels)
		},
	}


	return cmd
}

func listCalendarLabels(globalFlags *types.GlobalFlags, flags *listCalendarLabelsFlags, cmd *cobra.Command, args []string) error {

res, err := maintenance.Maintenance(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

