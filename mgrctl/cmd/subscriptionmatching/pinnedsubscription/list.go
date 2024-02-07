package pinnedsubscription

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/subscriptionmatching/pinnedsubscription"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all PinnedSubscriptions",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, list)
		},
	}


	return cmd
}

func list(globalFlags *types.GlobalFlags, flags *listFlags, cmd *cobra.Command, args []string) error {

res, err := pinnedsubscription.Pinnedsubscription(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

