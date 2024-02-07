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

type setTimeZoneFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Login                 string
	Tzid                  int
}

func setTimeZoneCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setTimeZone",
		Short: "Set a user's timezone.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setTimeZoneFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setTimeZone)
		},
	}

	cmd.Flags().String("Login", "", "User's login name.")
	cmd.Flags().String("Tzid", "", "Timezone ID. (from listTimeZones)")

	return cmd
}

func setTimeZone(globalFlags *types.GlobalFlags, flags *setTimeZoneFlags, cmd *cobra.Command, args []string) error {

	res, err := locale.Locale(&flags.ConnectionDetails, flags.Login, flags.Tzid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
