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

type listUngroupedSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listUngroupedSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listUngroupedSystems",
		Short: "List systems that are not associated with any system groups.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listUngroupedSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listUngroupedSystems)
		},
	}


	return cmd
}

func listUngroupedSystems(globalFlags *types.GlobalFlags, flags *listUngroupedSystemsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

