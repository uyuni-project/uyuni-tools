package errata

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/errata"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type applicableToChannelsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	AdvisoryName          string
}

func applicableToChannelsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "applicableToChannels",
		Short: "Returns a list of channels applicable to the errata with the given advisory name.
 For those errata that are present in both vendor and user organizations under the same advisory name,
 this method retrieves the list of channels applicable of both of them.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags applicableToChannelsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, applicableToChannels)
		},
	}

	cmd.Flags().String("AdvisoryName", "", "")

	return cmd
}

func applicableToChannels(globalFlags *types.GlobalFlags, flags *applicableToChannelsFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.AdvisoryName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

