package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listUserSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login                 string
}

func listUserSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listUserSystems",
		Short: "List systems for a given user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listUserSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listUserSystems)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")

	return cmd
}

func listUserSystems(globalFlags *types.GlobalFlags, flags *listUserSystemsFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Login)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
