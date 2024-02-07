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

type applyChannelStateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sids          []int
}

func applyChannelStateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "applyChannelState",
		Short: "Refresh pillar data and then schedule channels state on the provided systems",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags applyChannelStateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, applyChannelState)
		},
	}

	cmd.Flags().String("Sids", "", "$desc")

	return cmd
}

func applyChannelState(globalFlags *types.GlobalFlags, flags *applyChannelStateFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.Sids)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

