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

type createChainFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChainLabel            string
}

func createChainCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createChain",
		Short: "Create an Action Chain.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createChainFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createChain)
		},
	}

	cmd.Flags().String("ChainLabel", "", "Label of the chain")

	return cmd
}

func createChain(globalFlags *types.GlobalFlags, flags *createChainFlags, cmd *cobra.Command, args []string) error {

	res, err := actionchain.Actionchain(&flags.ConnectionDetails, flags.ChainLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
