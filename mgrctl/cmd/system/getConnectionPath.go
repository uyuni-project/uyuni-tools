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

type getConnectionPathFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func getConnectionPathCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getConnectionPath",
		Short: "Get the list of proxies that the given system connects
 through in order to reach the server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getConnectionPathFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getConnectionPath)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getConnectionPath(globalFlags *types.GlobalFlags, flags *getConnectionPathFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

