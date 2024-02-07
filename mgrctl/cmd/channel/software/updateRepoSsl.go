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

type updateRepoSslFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Id          int
	SslCaCert          string
	SslCliCert          string
	SslCliKey          string
	Label          string
}

func updateRepoSslCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateRepoSsl",
		Short: "Updates repository SSL certificates",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateRepoSslFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, updateRepoSsl)
		},
	}

	cmd.Flags().String("Id", "", "repository ID")
	cmd.Flags().String("SslCaCert", "", "SSL CA cert description")
	cmd.Flags().String("SslCliCert", "", "SSL Client cert description")
	cmd.Flags().String("SslCliKey", "", "SSL Client key description")
	cmd.Flags().String("Label", "", "repository label")

	return cmd
}

func updateRepoSsl(globalFlags *types.GlobalFlags, flags *updateRepoSslFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.Id, flags.SslCaCert, flags.SslCliCert, flags.SslCliKey, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

