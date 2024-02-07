package admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/admin"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type createFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Description           string
	Host                  string
	Port                  int
	Username              string
	Password              string
	Key                   string
	KeyPassword           string
	BastionHost           string
	BastionPort           int
	BastionUsername       string
	BastionPassword       string
	BastionKey            string
	BastionKeyPassword    string
}

func createCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new ssh connection data to extract data from",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, create)
		},
	}

	cmd.Flags().String("Description", "", "")
	cmd.Flags().String("Host", "", "hostname or IP address to the instance, will fail if already in use.")
	cmd.Flags().String("Port", "", "")
	cmd.Flags().String("Username", "", "")
	cmd.Flags().String("Password", "", "")
	cmd.Flags().String("Key", "", "private key to use in authentication")
	cmd.Flags().String("KeyPassword", "", "")
	cmd.Flags().String("BastionHost", "", "hostname or IP address to a bastion host")
	cmd.Flags().String("BastionPort", "", "")
	cmd.Flags().String("BastionUsername", "", "")
	cmd.Flags().String("BastionPassword", "", "")
	cmd.Flags().String("BastionKey", "", "private key to use in bastion authentication")
	cmd.Flags().String("BastionKeyPassword", "", "")

	return cmd
}

func create(globalFlags *types.GlobalFlags, flags *createFlags, cmd *cobra.Command, args []string) error {

	res, err := admin.Admin(&flags.ConnectionDetails, flags.Description, flags.Host, flags.Port, flags.Username, flags.Password, flags.Key, flags.KeyPassword, flags.BastionHost, flags.BastionPort, flags.BastionUsername, flags.BastionPassword, flags.BastionKey, flags.BastionKeyPassword)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
