package access

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/access"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setOrgSharingFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	Access                string
}

func setOrgSharingCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setOrgSharing",
		Short: "Set organization sharing access control.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setOrgSharingFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setOrgSharing)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "label of the channel")
	cmd.Flags().String("Access", "", "Access (one of the                  following: 'public', 'private', or 'protected'")

	return cmd
}

func setOrgSharing(globalFlags *types.GlobalFlags, flags *setOrgSharingFlags, cmd *cobra.Command, args []string) error {

	res, err := access.Access(&flags.ConnectionDetails, flags.ChannelLabel, flags.Access)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
