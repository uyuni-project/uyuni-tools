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

type deleteChainFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChainLabel            string
}

func deleteChainCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteChain",
		Short: "Delete action chain by label.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteChainFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteChain)
		},
	}

	cmd.Flags().String("ChainLabel", "", "Label of the chain")

	return cmd
}

func deleteChain(globalFlags *types.GlobalFlags, flags *deleteChainFlags, cmd *cobra.Command, args []string) error {

	res, err := actionchain.Actionchain(&flags.ConnectionDetails, flags.ChainLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
