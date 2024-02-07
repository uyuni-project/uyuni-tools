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

type createCalendarWithUrlFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	Url          string
}

func createCalendarWithUrlCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createCalendarWithUrl",
		Short: "Create a new maintenance calendar",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createCalendarWithUrlFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createCalendarWithUrl)
		},
	}

	cmd.Flags().String("Label", "", "maintenance calendar label")
	cmd.Flags().String("Url", "", "download URL for ICal calendar data")

	return cmd
}

func createCalendarWithUrl(globalFlags *types.GlobalFlags, flags *createCalendarWithUrlFlags, cmd *cobra.Command, args []string) error {

res, err := maintenance.Maintenance(&flags.ConnectionDetails, flags.Label, flags.Url)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

