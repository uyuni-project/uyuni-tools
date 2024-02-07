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

type createFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login          string
	Password          string
	FirstName          string
	LastName          string
	Email          string
	Login          string
	UsePamAuth          int
}

func createCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, create)
		},
	}

	cmd.Flags().String("Login", "", "desired login name, will fail if already in use.")
	cmd.Flags().String("Password", "", "")
	cmd.Flags().String("FirstName", "", "")
	cmd.Flags().String("LastName", "", "")
	cmd.Flags().String("Email", "", "User's e-mail address.")
	cmd.Flags().String("Login", "", "desired login name, will fail if already in use.")
	cmd.Flags().String("UsePamAuth", "", "1 if you wish to use PAM authentication for this user, 0 otherwise.")

	return cmd
}

func create(globalFlags *types.GlobalFlags, flags *createFlags, cmd *cobra.Command, args []string) error {

res, err := user.User(&flags.ConnectionDetails, flags.Login, flags.Password, flags.FirstName, flags.LastName, flags.Email, flags.Login, flags.UsePamAuth)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

