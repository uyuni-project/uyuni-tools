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

type listLatestPackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func listLatestPackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listLatestPackages",
		Short: "Lists the packages with the latest version (including release and
 epoch) for the given channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listLatestPackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listLatestPackages)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "channel to query")

	return cmd
}

func listLatestPackages(globalFlags *types.GlobalFlags, flags *listLatestPackagesFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

