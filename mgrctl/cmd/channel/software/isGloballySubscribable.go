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

type isGloballySubscribableFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func isGloballySubscribableCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "isGloballySubscribable",
		Short: "Returns whether the channel is subscribable by any user
 in the organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags isGloballySubscribableFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, isGloballySubscribable)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "channel to query")

	return cmd
}

func isGloballySubscribable(globalFlags *types.GlobalFlags, flags *isGloballySubscribableFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

