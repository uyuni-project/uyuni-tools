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

type updateRepoUrlFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Id                    int
	Url                   string
	Label                 string
}

func updateRepoUrlCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateRepoUrl",
		Short: "Updates repository source URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateRepoUrlFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, updateRepoUrl)
		},
	}

	cmd.Flags().String("Id", "", "repository ID")
	cmd.Flags().String("Url", "", "new repository URL")
	cmd.Flags().String("Label", "", "repository label")

	return cmd
}

func updateRepoUrl(globalFlags *types.GlobalFlags, flags *updateRepoUrlFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.Id, flags.Url, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
