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

type listLocalesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listLocalesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listLocales",
		Short: "Returns a list of all understood locales. Can be
 used as input to setLocale.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listLocalesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listLocales)
		},
	}


	return cmd
}

func listLocales(globalFlags *types.GlobalFlags, flags *listLocalesFlags, cmd *cobra.Command, args []string) error {

res, err := locale.Locale(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

