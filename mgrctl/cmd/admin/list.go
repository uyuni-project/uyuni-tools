package admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/admin"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Returns a list of ssh connection data registered.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, list)
		},
	}

	return cmd
}

func list(globalFlags *types.GlobalFlags, flags *listFlags, cmd *cobra.Command, args []string) error {

	res, err := admin.Admin(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
