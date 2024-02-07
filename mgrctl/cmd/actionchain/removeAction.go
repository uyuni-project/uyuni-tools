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

type removeActionFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChainLabel          string
	ActionId          int
}

func removeActionCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeAction",
		Short: "Remove an action from an Action Chain.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeActionFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeAction)
		},
	}

	cmd.Flags().String("ChainLabel", "", "Label of the chain")
	cmd.Flags().String("ActionId", "", "Action ID")

	return cmd
}

func removeAction(globalFlags *types.GlobalFlags, flags *removeActionFlags, cmd *cobra.Command, args []string) error {

res, err := actionchain.Actionchain(&flags.ConnectionDetails, flags.ChainLabel, flags.ActionId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

