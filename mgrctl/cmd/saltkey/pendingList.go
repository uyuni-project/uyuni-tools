package saltkey

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/saltkey"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type pendingListFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func pendingListCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pendingList",
		Short: "List pending salt keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags pendingListFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, pendingList)
		},
	}

	return cmd
}

func pendingList(globalFlags *types.GlobalFlags, flags *pendingListFlags, cmd *cobra.Command, args []string) error {

	res, err := saltkey.Saltkey(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
