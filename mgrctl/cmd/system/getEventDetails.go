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

type getEventDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	Eid          int
}

func getEventDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getEventDetails",
		Short: "Returns the details of the event associated with the specified server and event.
             The event id must be a value returned by the system.getEventHistory API.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getEventDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getEventDetails)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Eid", "", "ID of the event")

	return cmd
}

func getEventDetails(globalFlags *types.GlobalFlags, flags *getEventDetailsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.Eid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

