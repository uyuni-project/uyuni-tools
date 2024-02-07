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

type setReadOnlyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login          string
	ReadOnly          bool
}

func setReadOnlyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setReadOnly",
		Short: "Sets whether the target user should have only read-only API access or
 standard full scale access.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setReadOnlyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setReadOnly)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")
	cmd.Flags().String("ReadOnly", "", "Sets whether the target user should have only read-only API access or standard full scale access.")

	return cmd
}

func setReadOnly(globalFlags *types.GlobalFlags, flags *setReadOnlyFlags, cmd *cobra.Command, args []string) error {

res, err := user.User(&flags.ConnectionDetails, flags.Login, flags.ReadOnly)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

