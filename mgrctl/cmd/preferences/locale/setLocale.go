package locale

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/preferences/locale"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setLocaleFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login                 string
	Locale                string
}

func setLocaleCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setLocale",
		Short: "Set a user's locale.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setLocaleFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setLocale)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")
	cmd.Flags().String("Locale", "", "Locale to set. (from listLocales)")

	return cmd
}

func setLocale(globalFlags *types.GlobalFlags, flags *setLocaleFlags, cmd *cobra.Command, args []string) error {

	res, err := locale.Locale(&flags.ConnectionDetails, flags.Login, flags.Locale)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
