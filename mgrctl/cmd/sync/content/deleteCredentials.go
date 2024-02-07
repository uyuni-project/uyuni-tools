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

type deleteCredentialsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Username          string
}

func deleteCredentialsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteCredentials",
		Short: "Delete organization credentials (mirror credentials) from #product().",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteCredentialsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteCredentials)
		},
	}

	cmd.Flags().String("Username", "", "Username of credentials to delete")

	return cmd
}

func deleteCredentials(globalFlags *types.GlobalFlags, flags *deleteCredentialsFlags, cmd *cobra.Command, args []string) error {

res, err := content.Content(&flags.ConnectionDetails, flags.Username)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

