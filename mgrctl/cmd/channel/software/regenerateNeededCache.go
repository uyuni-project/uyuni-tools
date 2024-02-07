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

type regenerateNeededCacheFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func regenerateNeededCacheCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "regenerateNeededCache",
		Short: "Completely clear and regenerate the needed Errata and Package
      cache for all systems subscribed to the specified channel.  This should
      be used only if you believe your cache is incorrect for all the systems
      in a given channel. This will schedule an asynchronous action to actually
      do the processing.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags regenerateNeededCacheFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, regenerateNeededCache)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "the label of the          channel")

	return cmd
}

func regenerateNeededCache(globalFlags *types.GlobalFlags, flags *regenerateNeededCacheFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

