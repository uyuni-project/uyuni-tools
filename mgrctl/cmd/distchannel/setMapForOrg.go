package distchannel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/distchannel"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setMapForOrgFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Os          string
	Release          string
	ArchName          string
	ChannelLabel          string
}

func setMapForOrgCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setMapForOrg",
		Short: "Sets, overrides (/removes if channelLabel empty)
 a distribution channel map within an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setMapForOrgFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setMapForOrg)
		},
	}

	cmd.Flags().String("Os", "", "")
	cmd.Flags().String("Release", "", "")
	cmd.Flags().String("ArchName", "", "")
	cmd.Flags().String("ChannelLabel", "", "")

	return cmd
}

func setMapForOrg(globalFlags *types.GlobalFlags, flags *setMapForOrgFlags, cmd *cobra.Command, args []string) error {

res, err := distchannel.Distchannel(&flags.ConnectionDetails, flags.Os, flags.Release, flags.ArchName, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

