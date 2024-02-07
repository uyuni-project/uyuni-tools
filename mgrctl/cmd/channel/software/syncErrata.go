package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type syncErrataFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func syncErrataCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "syncErrata",
		Short: "If you have synced a new channel then patches
 will have been updated with the packages that are in the newly
 synced channel. A cloned erratum will not have been automatically updated
 however. If you cloned a channel that includes those cloned errata and
 should include the new packages, they will not be included when they
 should. This method updates all the errata in the given cloned channel
 with packages that have recently been added, and ensures that all the
 packages you expect are in the channel. It also updates cloned errata
 attributes like advisoryStatus.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags syncErrataFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, syncErrata)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "channel to update")

	return cmd
}

func syncErrata(globalFlags *types.GlobalFlags, flags *syncErrataFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

