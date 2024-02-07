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

type removeRepoFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Id          long
	Label          string
}

func removeRepoCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeRepo",
		Short: "Removes a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeRepoFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeRepo)
		},
	}

	cmd.Flags().String("Id", "", "ID of repo to be removed")
	cmd.Flags().String("Label", "", "label of repo to be removed")

	return cmd
}

func removeRepo(globalFlags *types.GlobalFlags, flags *removeRepoFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.Id, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

