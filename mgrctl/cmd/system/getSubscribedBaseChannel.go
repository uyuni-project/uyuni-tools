package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getSubscribedBaseChannelFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func getSubscribedBaseChannelCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getSubscribedBaseChannel",
		Short: "Provides the base channel of a given system",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getSubscribedBaseChannelFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getSubscribedBaseChannel)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getSubscribedBaseChannel(globalFlags *types.GlobalFlags, flags *getSubscribedBaseChannelFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
