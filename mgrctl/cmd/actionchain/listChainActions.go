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

type listChainActionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChainLabel          string
}

func listChainActionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listChainActions",
		Short: "List all actions in the particular Action Chain.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listChainActionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listChainActions)
		},
	}

	cmd.Flags().String("ChainLabel", "", "Label of the chain")

	return cmd
}

func listChainActions(globalFlags *types.GlobalFlags, flags *listChainActionsFlags, cmd *cobra.Command, args []string) error {

res, err := actionchain.Actionchain(&flags.ConnectionDetails, flags.ChainLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

