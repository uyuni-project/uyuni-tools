package systemgroup

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/systemgroup"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type subscribeConfigChannelFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemGroupName       string
	ConfigChannelLabels   []string
}

func subscribeConfigChannelCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscribeConfigChannel",
		Short: "Subscribe given config channels to a system group",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags subscribeConfigChannelFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, subscribeConfigChannel)
		},
	}

	cmd.Flags().String("SystemGroupName", "", "")
	cmd.Flags().String("ConfigChannelLabels", "", "$desc")

	return cmd
}

func subscribeConfigChannel(globalFlags *types.GlobalFlags, flags *subscribeConfigChannelFlags, cmd *cobra.Command, args []string) error {

	res, err := systemgroup.Systemgroup(&flags.ConnectionDetails, flags.SystemGroupName, flags.ConfigChannelLabels)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
