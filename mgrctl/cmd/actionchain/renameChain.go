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

type renameChainFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	PreviousLabel          string
	NewLabel          string
}

func renameChainCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renameChain",
		Short: "Rename an Action Chain.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags renameChainFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, renameChain)
		},
	}

	cmd.Flags().String("PreviousLabel", "", "Previous chain label")
	cmd.Flags().String("NewLabel", "", "New chain label")

	return cmd
}

func renameChain(globalFlags *types.GlobalFlags, flags *renameChainFlags, cmd *cobra.Command, args []string) error {

res, err := actionchain.Actionchain(&flags.ConnectionDetails, flags.PreviousLabel, flags.NewLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

