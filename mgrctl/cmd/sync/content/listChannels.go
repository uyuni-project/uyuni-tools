package content

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/sync/content"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listChannels",
		Short: "List all accessible channels.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listChannels)
		},
	}

	return cmd
}

func listChannels(globalFlags *types.GlobalFlags, flags *listChannelsFlags, cmd *cobra.Command, args []string) error {

	res, err := content.Content(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
