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

type downloadSystemIdFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func downloadSystemIdCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "downloadSystemId",
		Short: "Get the system ID file for a given server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags downloadSystemIdFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, downloadSystemId)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func downloadSystemId(globalFlags *types.GlobalFlags, flags *downloadSystemIdFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

