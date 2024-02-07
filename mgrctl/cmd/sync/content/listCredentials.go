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

type listCredentialsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listCredentialsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listCredentials",
		Short: "List organization credentials (mirror credentials) available in
             #product().",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listCredentialsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listCredentials)
		},
	}


	return cmd
}

func listCredentials(globalFlags *types.GlobalFlags, flags *listCredentialsFlags, cmd *cobra.Command, args []string) error {

res, err := content.Content(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

