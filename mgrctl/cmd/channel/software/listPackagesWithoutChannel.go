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

type listPackagesWithoutChannelFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listPackagesWithoutChannelCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listPackagesWithoutChannel",
		Short: "Lists all packages that are not associated with a channel.  Typically
          these are custom packages.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listPackagesWithoutChannelFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listPackagesWithoutChannel)
		},
	}


	return cmd
}

func listPackagesWithoutChannel(globalFlags *types.GlobalFlags, flags *listPackagesWithoutChannelFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

