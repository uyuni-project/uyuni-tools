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

type listPackagesFromChannelFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	ChannelLabel          string
}

func listPackagesFromChannelCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listPackagesFromChannel",
		Short: "Provides a list of packages installed on a system that are also
          contained in the given channel.  The installed package list did not
          include arch information before RHEL 5, so it is arch unaware.  RHEL 5
          systems do upload the arch information, and thus are arch aware.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listPackagesFromChannelFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listPackagesFromChannel)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("ChannelLabel", "", "")

	return cmd
}

func listPackagesFromChannel(globalFlags *types.GlobalFlags, flags *listPackagesFromChannelFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

