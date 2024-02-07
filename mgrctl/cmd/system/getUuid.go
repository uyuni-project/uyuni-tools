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

type getUuidFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func getUuidCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getUuid",
		Short: "Get the UUID from the given system ID.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getUuidFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getUuid)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getUuid(globalFlags *types.GlobalFlags, flags *getUuidFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
