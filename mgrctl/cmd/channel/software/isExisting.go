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

type isExistingFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func isExistingCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "isExisting",
		Short: "Returns whether is existing",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags isExistingFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, isExisting)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "label of the channel")

	return cmd
}

func isExisting(globalFlags *types.GlobalFlags, flags *isExistingFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

