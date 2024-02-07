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

type listDuplicatesByMacFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listDuplicatesByMacCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listDuplicatesByMac",
		Short: "List duplicate systems by Mac Address.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listDuplicatesByMacFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listDuplicatesByMac)
		},
	}

	return cmd
}

func listDuplicatesByMac(globalFlags *types.GlobalFlags, flags *listDuplicatesByMacFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
