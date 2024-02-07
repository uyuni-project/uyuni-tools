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

type disableFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login                 string
}

func disableCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable a user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags disableFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, disable)
		},
	}

	cmd.Flags().String("Login", "", "User login name to disable.")

	return cmd
}

func disable(globalFlags *types.GlobalFlags, flags *disableFlags, cmd *cobra.Command, args []string) error {

	res, err := user.User(&flags.ConnectionDetails, flags.Login)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
