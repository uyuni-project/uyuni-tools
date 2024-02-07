package auth

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/auth"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type logoutFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func logoutCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout the user with the given session key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags logoutFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, logout)
		},
	}


	return cmd
}

func logout(globalFlags *types.GlobalFlags, flags *logoutFlags, cmd *cobra.Command, args []string) error {

res, err := auth.Auth(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

