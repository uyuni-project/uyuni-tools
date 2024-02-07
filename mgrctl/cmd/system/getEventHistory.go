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

type getEventHistoryFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	EarliestDate          $date
	Offset          int
	Limit          int
}

func getEventHistoryCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getEventHistory",
		Short: "Returns a list of history items associated with the system happened after the specified date.
             The list is paged and ordered from newest to oldest.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getEventHistoryFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getEventHistory)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("EarliestDate", "", "")
	cmd.Flags().String("Offset", "", "Number of results to skip")
	cmd.Flags().String("Limit", "", "Maximum number of results")

	return cmd
}

func getEventHistory(globalFlags *types.GlobalFlags, flags *getEventHistoryFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.EarliestDate, flags.Offset, flags.Limit)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

