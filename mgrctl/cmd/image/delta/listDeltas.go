package delta

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/image/delta"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listDeltasFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listDeltasCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listDeltas",
		Short: "List available DeltaImages",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listDeltasFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listDeltas)
		},
	}


	return cmd
}

func listDeltas(globalFlags *types.GlobalFlags, flags *listDeltasFlags, cmd *cobra.Command, args []string) error {

res, err := delta.Delta(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

