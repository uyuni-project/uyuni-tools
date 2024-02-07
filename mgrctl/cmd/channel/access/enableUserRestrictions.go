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

type enableUserRestrictionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func enableUserRestrictionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enableUserRestrictions",
		Short: "Enable user restrictions for the given channel. If enabled, only
 selected users within the organization may subscribe to the channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags enableUserRestrictionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, enableUserRestrictions)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "label of the channel")

	return cmd
}

func enableUserRestrictions(globalFlags *types.GlobalFlags, flags *enableUserRestrictionsFlags, cmd *cobra.Command, args []string) error {

res, err := access.Access(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

