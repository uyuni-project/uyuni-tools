package user

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/user"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type usePamAuthenticationFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login          string
}

func usePamAuthenticationCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "usePamAuthentication",
		Short: "Toggles whether or not a user uses PAM authentication or
 basic #product() authentication.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags usePamAuthenticationFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, usePamAuthentication)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")

	return cmd
}

func usePamAuthentication(globalFlags *types.GlobalFlags, flags *usePamAuthenticationFlags, cmd *cobra.Command, args []string) error {

res, err := user.User(&flags.ConnectionDetails, flags.Login)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

