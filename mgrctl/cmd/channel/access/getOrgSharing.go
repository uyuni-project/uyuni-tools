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

type getOrgSharingFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func getOrgSharingCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getOrgSharing",
		Short: "Get organization sharing access control.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getOrgSharingFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getOrgSharing)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "label of the channel")

	return cmd
}

func getOrgSharing(globalFlags *types.GlobalFlags, flags *getOrgSharingFlags, cmd *cobra.Command, args []string) error {

res, err := access.Access(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

