package kickstart

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type findKickstartForIpFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	IpAddress             string
}

func findKickstartForIpCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "findKickstartForIp",
		Short: "Find an associated kickstart for a given ip address.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags findKickstartForIpFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, findKickstartForIp)
		},
	}

	cmd.Flags().String("IpAddress", "", "The ip address to search for (i.e. 192.168.0.1)")

	return cmd
}

func findKickstartForIp(globalFlags *types.GlobalFlags, flags *findKickstartForIpFlags, cmd *cobra.Command, args []string) error {

	res, err := kickstart.Kickstart(&flags.ConnectionDetails, flags.IpAddress)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
