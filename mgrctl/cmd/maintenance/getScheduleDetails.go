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

type getScheduleDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name                  string
}

func getScheduleDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getScheduleDetails",
		Short: "Lookup a specific maintenance schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getScheduleDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getScheduleDetails)
		},
	}

	cmd.Flags().String("Name", "", "maintenance Schedule Name")

	return cmd
}

func getScheduleDetails(globalFlags *types.GlobalFlags, flags *getScheduleDetailsFlags, cmd *cobra.Command, args []string) error {

	res, err := maintenance.Maintenance(&flags.ConnectionDetails, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
