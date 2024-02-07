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

type enableFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login                 string
}

func enableCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable a user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags enableFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, enable)
		},
	}

	cmd.Flags().String("Login", "", "User login name to enable.")

	return cmd
}

func enable(globalFlags *types.GlobalFlags, flags *enableFlags, cmd *cobra.Command, args []string) error {

	res, err := user.User(&flags.ConnectionDetails, flags.Login)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
