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

type listChannelReposFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
}

func listChannelReposCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listChannelRepos",
		Short: "Lists associated repos with the given channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listChannelReposFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listChannelRepos)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "channel label")

	return cmd
}

func listChannelRepos(globalFlags *types.GlobalFlags, flags *listChannelReposFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
