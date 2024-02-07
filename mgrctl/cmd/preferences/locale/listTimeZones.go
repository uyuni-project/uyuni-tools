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

type listTimeZonesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listTimeZonesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listTimeZones",
		Short: "Returns a list of all understood timezones. Results can be
 used as input to setTimeZone.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listTimeZonesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listTimeZones)
		},
	}


	return cmd
}

func listTimeZones(globalFlags *types.GlobalFlags, flags *listTimeZonesFlags, cmd *cobra.Command, args []string) error {

res, err := locale.Locale(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

