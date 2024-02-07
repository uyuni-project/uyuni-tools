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

type setContactDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	MaintainerName          string
	MaintainerEmail          string
	MaintainerPhone          string
	SupportPolicy          string
}

func setContactDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setContactDetails",
		Short: "Set contact/support information for given channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setContactDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setContactDetails)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "label of the channel")
	cmd.Flags().String("MaintainerName", "", "name of the channel maintainer")
	cmd.Flags().String("MaintainerEmail", "", "email of the channel maintainer")
	cmd.Flags().String("MaintainerPhone", "", "phone number of the channel maintainer")
	cmd.Flags().String("SupportPolicy", "", "channel support policy")

	return cmd
}

func setContactDetails(globalFlags *types.GlobalFlags, flags *setContactDetailsFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel, flags.MaintainerName, flags.MaintainerEmail, flags.MaintainerPhone, flags.SupportPolicy)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

