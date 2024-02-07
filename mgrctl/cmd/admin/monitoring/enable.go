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

type enableFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func enableCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable monitoring.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags enableFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, enable)
		},
	}


	return cmd
}

func enable(globalFlags *types.GlobalFlags, flags *enableFlags, cmd *cobra.Command, args []string) error {

res, err := monitoring.Monitoring(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

