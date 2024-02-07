package monitoring

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/admin/monitoring"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getStatusFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func getStatusCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getStatus",
		Short: "Get the status of each Prometheus exporter.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getStatusFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getStatus)
		},
	}

	return cmd
}

func getStatus(globalFlags *types.GlobalFlags, flags *getStatusFlags, cmd *cobra.Command, args []string) error {

	res, err := monitoring.Monitoring(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
