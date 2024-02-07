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

type setUserManageableFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	Login          string
	Value          bool
}

func setUserManageableCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setUserManageable",
		Short: "Set the manageable flag for a given channel and user.
 If value is set to 'true', this method will give the user
 manage permissions to the channel. Otherwise, that privilege is revoked.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setUserManageableFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setUserManageable)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "label of the channel")
	cmd.Flags().String("Login", "", "login of the target user")
	cmd.Flags().String("Value", "", "value of the flag to set")

	return cmd
}

func setUserManageable(globalFlags *types.GlobalFlags, flags *setUserManageableFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel, flags.Login, flags.Value)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

