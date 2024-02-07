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

type associateRepoFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	RepoLabel             string
}

func associateRepoCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "associateRepo",
		Short: "Associates a repository with a channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags associateRepoFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, associateRepo)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "channel label")
	cmd.Flags().String("RepoLabel", "", "repository label")

	return cmd
}

func associateRepo(globalFlags *types.GlobalFlags, flags *associateRepoFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel, flags.RepoLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
