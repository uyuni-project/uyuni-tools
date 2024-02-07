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

type deleteFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Host                  string
}

func deleteCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Returns a list of ssh connection data registered.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, delete)
		},
	}

	cmd.Flags().String("Host", "", "hostname or IP address of the instance")

	return cmd
}

func delete(globalFlags *types.GlobalFlags, flags *deleteFlags, cmd *cobra.Command, args []string) error {

	res, err := admin.Admin(&flags.ConnectionDetails, flags.Host)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
