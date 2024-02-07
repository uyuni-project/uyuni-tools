package activationkey

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/activationkey"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listActivatedSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Key                   string
}

func listActivatedSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listActivatedSystems",
		Short: "List the systems activated with the key provided.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listActivatedSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listActivatedSystems)
		},
	}

	cmd.Flags().String("Key", "", "")

	return cmd
}

func listActivatedSystems(globalFlags *types.GlobalFlags, flags *listActivatedSystemsFlags, cmd *cobra.Command, args []string) error {

	res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Key)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
