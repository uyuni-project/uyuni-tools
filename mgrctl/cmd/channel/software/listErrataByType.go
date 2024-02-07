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

type listErrataByTypeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	AdvisoryType          string
}

func listErrataByTypeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listErrataByType",
		Short: "List the errata of a specific type that are applicable to a channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listErrataByTypeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listErrataByType)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "channel to query")
	cmd.Flags().String("AdvisoryType", "", "type of advisory (one of of the following: 'Security Advisory', 'Product Enhancement Advisory', 'Bug Fix Advisory'")

	return cmd
}

func listErrataByType(globalFlags *types.GlobalFlags, flags *listErrataByTypeFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel, flags.AdvisoryType)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
