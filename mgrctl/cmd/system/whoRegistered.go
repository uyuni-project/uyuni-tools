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

type whoRegisteredFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func whoRegisteredCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoRegistered",
		Short: "Returns information about the user who registered the system",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags whoRegisteredFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, whoRegistered)
		},
	}

	cmd.Flags().String("Sid", "", "Id of the system in question")

	return cmd
}

func whoRegistered(globalFlags *types.GlobalFlags, flags *whoRegisteredFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
