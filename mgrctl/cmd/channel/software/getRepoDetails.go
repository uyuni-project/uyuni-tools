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

type getRepoDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	RepoLabel          string
	Id          int
}

func getRepoDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getRepoDetails",
		Short: "Returns details of the given repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getRepoDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getRepoDetails)
		},
	}

	cmd.Flags().String("RepoLabel", "", "repo to query")
	cmd.Flags().String("Id", "", "repository ID")

	return cmd
}

func getRepoDetails(globalFlags *types.GlobalFlags, flags *getRepoDetailsFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.RepoLabel, flags.Id)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

