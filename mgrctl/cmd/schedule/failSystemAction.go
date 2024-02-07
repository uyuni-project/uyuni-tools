package schedule

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/schedule"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type failSystemActionFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
	ActionId              int
	Message               string
}

func failSystemActionCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "failSystemAction",
		Short: "Fail specific event on specified system",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags failSystemActionFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, failSystemAction)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("ActionId", "", "")
	cmd.Flags().String("Message", "", "")

	return cmd
}

func failSystemAction(globalFlags *types.GlobalFlags, flags *failSystemActionFlags, cmd *cobra.Command, args []string) error {

	res, err := schedule.Schedule(&flags.ConnectionDetails, flags.Sid, flags.ActionId, flags.Message)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
