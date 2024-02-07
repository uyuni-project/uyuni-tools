package monitoring

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/monitoring"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listEndpointsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sids                  []int
}

func listEndpointsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listEndpoints",
		Short: "Get the list of monitoring endpoint details.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listEndpointsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listEndpoints)
		},
	}

	cmd.Flags().String("Sids", "", "$desc")

	return cmd
}

func listEndpoints(globalFlags *types.GlobalFlags, flags *listEndpointsFlags, cmd *cobra.Command, args []string) error {

	res, err := monitoring.Monitoring(&flags.ConnectionDetails, flags.Sids)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
