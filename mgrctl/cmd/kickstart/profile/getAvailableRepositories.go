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

type getAvailableRepositoriesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func getAvailableRepositoriesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getAvailableRepositories",
		Short: "Lists available OS repositories to associate with the provided
 kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getAvailableRepositoriesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getAvailableRepositories)
		},
	}

	cmd.Flags().String("KsLabel", "", "")

	return cmd
}

func getAvailableRepositories(globalFlags *types.GlobalFlags, flags *getAvailableRepositoriesFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

