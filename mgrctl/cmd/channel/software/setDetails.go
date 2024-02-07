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

type setDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	ChannelId             int
}

func setDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setDetails",
		Short: "Allows to modify channel attributes",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setDetails)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "channel label")
	cmd.Flags().String("ChannelId", "", "channel id")

	return cmd
}

func setDetails(globalFlags *types.GlobalFlags, flags *setDetailsFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel, flags.ChannelId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
