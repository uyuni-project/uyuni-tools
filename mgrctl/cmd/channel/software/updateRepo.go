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

type updateRepoFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Id          int
	Label          string
	Url          string
}

func updateRepoCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateRepo",
		Short: "Updates a ContentSource (repo)",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateRepoFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, updateRepo)
		},
	}

	cmd.Flags().String("Id", "", "repository ID")
	cmd.Flags().String("Label", "", "new repository label")
	cmd.Flags().String("Url", "", "new repository URL")

	return cmd
}

func updateRepo(globalFlags *types.GlobalFlags, flags *updateRepoFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.Id, flags.Label, flags.Url)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

