package tree

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/tree"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func listCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the available kickstartable trees for the given channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, list)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "Label of channel to search.")

	return cmd
}

func list(globalFlags *types.GlobalFlags, flags *listFlags, cmd *cobra.Command, args []string) error {

res, err := tree.Tree(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

