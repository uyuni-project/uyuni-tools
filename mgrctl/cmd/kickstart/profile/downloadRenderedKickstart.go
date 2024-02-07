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

type downloadRenderedKickstartFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel               string
}

func downloadRenderedKickstartCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "downloadRenderedKickstart",
		Short: "Downloads the Cobbler-rendered Kickstart file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags downloadRenderedKickstartFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, downloadRenderedKickstart)
		},
	}

	cmd.Flags().String("KsLabel", "", "The label of the kickstart to download.")

	return cmd
}

func downloadRenderedKickstart(globalFlags *types.GlobalFlags, flags *downloadRenderedKickstartFlags, cmd *cobra.Command, args []string) error {

	res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
