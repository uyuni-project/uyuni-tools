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

type listAllPackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	StartDate          $type
	EndDate          $type
}

func listAllPackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAllPackages",
		Short: "Lists all packages in the channel, regardless of package version,
 between the given dates.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAllPackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAllPackages)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "channel to query")
	cmd.Flags().String("StartDate", "", "")
	cmd.Flags().String("EndDate", "", "")

	return cmd
}

func listAllPackages(globalFlags *types.GlobalFlags, flags *listAllPackagesFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel, flags.StartDate, flags.EndDate)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

