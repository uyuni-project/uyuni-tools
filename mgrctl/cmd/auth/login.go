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

type loginFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Username          string
	Password          string
	Duration          int
}

func loginCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login using a username and password. Returns the session key
 used by most other API methods.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags loginFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, login)
		},
	}

	cmd.Flags().String("Username", "", "")
	cmd.Flags().String("Password", "", "")
	cmd.Flags().String("Duration", "", "Length of session.")

	return cmd
}

func login(globalFlags *types.GlobalFlags, flags *loginFlags, cmd *cobra.Command, args []string) error {

res, err := auth.Auth(&flags.ConnectionDetails, flags.Username, flags.Password, flags.Duration)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

