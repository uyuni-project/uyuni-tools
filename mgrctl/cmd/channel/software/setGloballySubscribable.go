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

type setGloballySubscribableFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	Value                 bool
}

func setGloballySubscribableCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setGloballySubscribable",
		Short: "Set globally subscribable attribute for given channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setGloballySubscribableFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setGloballySubscribable)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "label of the channel")
	cmd.Flags().String("Value", "", "true if the channel is to be          globally subscribable. False otherwise.")

	return cmd
}

func setGloballySubscribable(globalFlags *types.GlobalFlags, flags *setGloballySubscribableFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel, flags.Value)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
