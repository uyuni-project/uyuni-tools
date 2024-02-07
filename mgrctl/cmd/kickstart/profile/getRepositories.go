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

type getRepositoriesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func getRepositoriesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getRepositories",
		Short: "Lists all OS repositories associated with provided kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getRepositoriesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getRepositories)
		},
	}

	cmd.Flags().String("KsLabel", "", "")

	return cmd
}

func getRepositories(globalFlags *types.GlobalFlags, flags *getRepositoriesFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

