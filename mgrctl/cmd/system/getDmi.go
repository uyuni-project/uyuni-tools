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

type getDmiFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func getDmiCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getDmi",
		Short: "Gets the DMI information of a system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getDmiFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getDmi)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getDmi(globalFlags *types.GlobalFlags, flags *getDmiFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
