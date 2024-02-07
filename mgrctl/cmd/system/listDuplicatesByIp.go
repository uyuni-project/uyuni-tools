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

type listDuplicatesByIpFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listDuplicatesByIpCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listDuplicatesByIp",
		Short: "List duplicate systems by IP Address.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listDuplicatesByIpFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listDuplicatesByIp)
		},
	}

	return cmd
}

func listDuplicatesByIp(globalFlags *types.GlobalFlags, flags *listDuplicatesByIpFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
