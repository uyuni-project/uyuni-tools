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

type disassociateRepoFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	RepoLabel             string
}

func disassociateRepoCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disassociateRepo",
		Short: "Disassociates a repository from a channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags disassociateRepoFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, disassociateRepo)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "channel label")
	cmd.Flags().String("RepoLabel", "", "repository label")

	return cmd
}

func disassociateRepo(globalFlags *types.GlobalFlags, flags *disassociateRepoFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel, flags.RepoLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
