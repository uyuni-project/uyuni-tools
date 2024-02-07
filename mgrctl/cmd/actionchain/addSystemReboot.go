package actionchain

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/actionchain"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type addSystemRebootFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	ChainLabel          string
}

func addSystemRebootCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addSystemReboot",
		Short: "Add system reboot to an Action Chain.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addSystemRebootFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addSystemReboot)
		},
	}

	cmd.Flags().String("Sid", "", "System ID")
	cmd.Flags().String("ChainLabel", "", "Label of the chain")

	return cmd
}

func addSystemReboot(globalFlags *types.GlobalFlags, flags *addSystemRebootFlags, cmd *cobra.Command, args []string) error {

res, err := actionchain.Actionchain(&flags.ConnectionDetails, flags.Sid, flags.ChainLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

