package content

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/sync/content"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type synchronizeRepositoriesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	MirrorUrl          string
}

func synchronizeRepositoriesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "synchronizeRepositories",
		Short: "Synchronize repositories between the Customer Center
             and the #product() database.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags synchronizeRepositoriesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, synchronizeRepositories)
		},
	}

	cmd.Flags().String("MirrorUrl", "", "Optional mirror url or null")

	return cmd
}

func synchronizeRepositories(globalFlags *types.GlobalFlags, flags *synchronizeRepositoriesFlags, cmd *cobra.Command, args []string) error {

res, err := content.Content(&flags.ConnectionDetails, flags.MirrorUrl)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

