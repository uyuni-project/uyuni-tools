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

type getCalendarDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label                 string
}

func getCalendarDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getCalendarDetails",
		Short: "Lookup a specific maintenance schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getCalendarDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getCalendarDetails)
		},
	}

	cmd.Flags().String("Label", "", "maintenance calendar label")

	return cmd
}

func getCalendarDetails(globalFlags *types.GlobalFlags, flags *getCalendarDetailsFlags, cmd *cobra.Command, args []string) error {

	res, err := maintenance.Maintenance(&flags.ConnectionDetails, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
