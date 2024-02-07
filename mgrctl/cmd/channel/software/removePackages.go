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

type removePackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func removePackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removePackages",
		Short: "Removes a given list of packages from the given channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removePackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removePackages)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "target channel")

	return cmd
}

func removePackages(globalFlags *types.GlobalFlags, flags *removePackagesFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
