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

type setUserSubscribableFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	Login          string
	Value          bool
}

func setUserSubscribableCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setUserSubscribable",
		Short: "Set the subscribable flag for a given channel and user.
 If value is set to 'true', this method will give the user
 subscribe permissions to the channel. Otherwise, that privilege is revoked.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setUserSubscribableFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setUserSubscribable)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "label of the channel")
	cmd.Flags().String("Login", "", "login of the target user")
	cmd.Flags().String("Value", "", "value of the flag to set")

	return cmd
}

func setUserSubscribable(globalFlags *types.GlobalFlags, flags *setUserSubscribableFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel, flags.Login, flags.Value)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

