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

type listDuplicatesByHostnameFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listDuplicatesByHostnameCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listDuplicatesByHostname",
		Short: "List duplicate systems by Hostname.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listDuplicatesByHostnameFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listDuplicatesByHostname)
		},
	}


	return cmd
}

func listDuplicatesByHostname(globalFlags *types.GlobalFlags, flags *listDuplicatesByHostnameFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

