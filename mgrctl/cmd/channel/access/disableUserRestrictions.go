package access

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/access"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type disableUserRestrictionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func disableUserRestrictionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disableUserRestrictions",
		Short: "Disable user restrictions for the given channel.  If disabled,
 all users within the organization may subscribe to the channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags disableUserRestrictionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, disableUserRestrictions)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "label of the channel")

	return cmd
}

func disableUserRestrictions(globalFlags *types.GlobalFlags, flags *disableUserRestrictionsFlags, cmd *cobra.Command, args []string) error {

res, err := access.Access(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

