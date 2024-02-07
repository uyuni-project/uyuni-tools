package user

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/api/user"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setErrataNotificationsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login                 string
	Value                 bool
}

func setErrataNotificationsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setErrataNotifications",
		Short: "Enables/disables errata mail notifications for a specific user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setErrataNotificationsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setErrataNotifications)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")
	cmd.Flags().String("Value", "", "True for enabling errata notifications, False for disabling")

	return cmd
}

func setErrataNotifications(globalFlags *types.GlobalFlags, flags *setErrataNotificationsFlags, cmd *cobra.Command, args []string) error {

	res, err := user.User(&flags.ConnectionDetails, flags.Login, flags.Value)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
