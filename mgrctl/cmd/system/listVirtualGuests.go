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

type listVirtualGuestsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func listVirtualGuestsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listVirtualGuests",
		Short: "Lists the virtual guests for a given virtual host",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listVirtualGuestsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listVirtualGuests)
		},
	}

	cmd.Flags().String("Sid", "", "the virtual host's id")

	return cmd
}

func listVirtualGuests(globalFlags *types.GlobalFlags, flags *listVirtualGuestsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

