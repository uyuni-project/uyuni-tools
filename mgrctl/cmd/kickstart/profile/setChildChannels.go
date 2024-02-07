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

type setChildChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel               string
}

func setChildChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setChildChannels",
		Short: "Set the child channels for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setChildChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setChildChannels)
		},
	}

	cmd.Flags().String("KsLabel", "", "Label of kickstart profile to be changed.")

	return cmd
}

func setChildChannels(globalFlags *types.GlobalFlags, flags *setChildChannelsFlags, cmd *cobra.Command, args []string) error {

	res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
