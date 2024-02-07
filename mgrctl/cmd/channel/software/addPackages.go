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

type addPackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func addPackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addPackages",
		Short: "Adds a given list of packages to the given channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addPackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addPackages)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "target channel")

	return cmd
}

func addPackages(globalFlags *types.GlobalFlags, flags *addPackagesFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

