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

type listChainsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listChainsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listChains",
		Short: "List currently available action chains.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listChainsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listChains)
		},
	}

	return cmd
}

func listChains(globalFlags *types.GlobalFlags, flags *listChainsFlags, cmd *cobra.Command, args []string) error {

	res, err := actionchain.Actionchain(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
