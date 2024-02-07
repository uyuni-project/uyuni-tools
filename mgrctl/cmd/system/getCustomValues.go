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

type getCustomValuesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func getCustomValuesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getCustomValues",
		Short: "Get the custom data values defined for the server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getCustomValuesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getCustomValues)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getCustomValues(globalFlags *types.GlobalFlags, flags *getCustomValuesFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

