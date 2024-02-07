package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listChildrenFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func listChildrenCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listChildren",
		Short: "List the children of a channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listChildrenFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listChildren)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "the label of the channel")

	return cmd
}

func listChildren(globalFlags *types.GlobalFlags, flags *listChildrenFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
