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

type isUserSubscribableFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	Login                 string
}

func isUserSubscribableCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "isUserSubscribable",
		Short: "Returns whether the channel may be subscribed to by the given user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags isUserSubscribableFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, isUserSubscribable)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "label of the channel")
	cmd.Flags().String("Login", "", "login of the target user")

	return cmd
}

func isUserSubscribable(globalFlags *types.GlobalFlags, flags *isUserSubscribableFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel, flags.Login)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
