package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setLocaleFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel               string
	Locale                string
}

func setLocaleCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setLocale",
		Short: "Sets the locale for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setLocaleFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setLocale)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")
	cmd.Flags().String("Locale", "", "the locale")

	return cmd
}

func setLocale(globalFlags *types.GlobalFlags, flags *setLocaleFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.KsLabel, flags.Locale)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
