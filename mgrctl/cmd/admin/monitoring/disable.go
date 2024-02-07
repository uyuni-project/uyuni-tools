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

type disableFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func disableCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable monitoring.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags disableFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, disable)
		},
	}


	return cmd
}

func disable(globalFlags *types.GlobalFlags, flags *disableFlags, cmd *cobra.Command, args []string) error {

res, err := monitoring.Monitoring(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

