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

type getMinionIdMapFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func getMinionIdMapCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getMinionIdMap",
		Short: "Return a map from Salt minion IDs to System IDs.
 Map entries are limited to systems that are visible by the current user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getMinionIdMapFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getMinionIdMap)
		},
	}


	return cmd
}

func getMinionIdMap(globalFlags *types.GlobalFlags, flags *getMinionIdMapFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

