package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type downloadKickstartFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel               string
	Host                  string
}

func downloadKickstartCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "downloadKickstart",
		Short: "Download the full contents of a kickstart file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags downloadKickstartFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, downloadKickstart)
		},
	}

	cmd.Flags().String("KsLabel", "", "The label of the kickstart to download.")
	cmd.Flags().String("Host", "", "The host to use when referring to the #product() server. Usually this should be the FQDN, but could be the ip address or shortname as well.")

	return cmd
}

func downloadKickstart(globalFlags *types.GlobalFlags, flags *downloadKickstartFlags, cmd *cobra.Command, args []string) error {

	res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel, flags.Host)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
