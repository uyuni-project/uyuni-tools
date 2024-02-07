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

type listErrataNeedingSyncFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func listErrataNeedingSyncCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listErrataNeedingSync",
		Short: "If you have synced a new channel then patches
 will have been updated with the packages that are in the newly
 synced channel. A cloned erratum will not have been automatically updated
 however. If you cloned a channel that includes those cloned errata and
 should include the new packages, they will not be included when they
 should. This method lists the errata that will be updated if you run the
 syncErrata method.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listErrataNeedingSyncFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listErrataNeedingSync)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "channel to update")

	return cmd
}

func listErrataNeedingSync(globalFlags *types.GlobalFlags, flags *listErrataNeedingSyncFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

