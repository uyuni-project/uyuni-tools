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

type setCustomValuesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	$param.getFlagName()          $param.getType()
}

func setCustomValuesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setCustomValues",
		Short: "Set custom values for the specified server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setCustomValuesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setCustomValues)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func setCustomValues(globalFlags *types.GlobalFlags, flags *setCustomValuesFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

