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

type regenerateYumCacheFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	Force                 bool
}

func regenerateYumCacheCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "regenerateYumCache",
		Short: "Regenerate yum cache for the specified channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags regenerateYumCacheFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, regenerateYumCache)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "the label of the          channel")
	cmd.Flags().String("Force", "", "force cache regeneration")

	return cmd
}

func regenerateYumCache(globalFlags *types.GlobalFlags, flags *regenerateYumCacheFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel, flags.Force)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
