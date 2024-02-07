package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setRepositoriesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel               string
	RepoLabels            []string
}

func setRepositoriesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setRepositories",
		Short: "Associates OS repository to a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setRepositoriesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setRepositories)
		},
	}

	cmd.Flags().String("KsLabel", "", "")
	cmd.Flags().String("RepoLabels", "", "$desc")

	return cmd
}

func setRepositories(globalFlags *types.GlobalFlags, flags *setRepositoriesFlags, cmd *cobra.Command, args []string) error {

	res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel, flags.RepoLabels)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
