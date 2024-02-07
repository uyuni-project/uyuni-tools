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

type listUserReposFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listUserReposCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listUserRepos",
		Short: "Returns a list of ContentSource (repos) that the user can see",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listUserReposFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listUserRepos)
		},
	}


	return cmd
}

func listUserRepos(globalFlags *types.GlobalFlags, flags *listUserReposFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

