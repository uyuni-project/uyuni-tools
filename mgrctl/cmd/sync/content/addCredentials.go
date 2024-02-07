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

type addCredentialsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Username              string
	Password              string
	Primary               bool
}

func addCredentialsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addCredentials",
		Short: "Add organization credentials (mirror credentials) to #product().",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addCredentialsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addCredentials)
		},
	}

	cmd.Flags().String("Username", "", "Organization credentials                                                  (Mirror credentials) username")
	cmd.Flags().String("Password", "", "Organization credentials                                                  (Mirror credentials) password")
	cmd.Flags().String("Primary", "", "Make this the primary credentials")

	return cmd
}

func addCredentials(globalFlags *types.GlobalFlags, flags *addCredentialsFlags, cmd *cobra.Command, args []string) error {

	res, err := content.Content(&flags.ConnectionDetails, flags.Username, flags.Password, flags.Primary)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
