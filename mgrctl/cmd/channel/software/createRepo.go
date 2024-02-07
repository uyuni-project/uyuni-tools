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

type createRepoFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label                 string
	Type                  string
	Url                   string
	Type                  string
	SslCaCert             string
	SslCliCert            string
	SslCliKey             string
	Type                  string
	SslCaCert             string
	SslCliCert            string
	SslCliKey             string
	HasSignedMetadata     bool
}

func createRepoCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createRepo",
		Short: "Creates a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createRepoFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createRepo)
		},
	}

	cmd.Flags().String("Label", "", "repository label")
	cmd.Flags().String("Type", "", "repository type (yum, uln...)")
	cmd.Flags().String("Url", "", "repository url")
	cmd.Flags().String("Type", "", "repository type (yum, uln...)")
	cmd.Flags().String("SslCaCert", "", "SSL CA cert description")
	cmd.Flags().String("SslCliCert", "", "SSL Client cert description")
	cmd.Flags().String("SslCliKey", "", "SSL Client key description")
	cmd.Flags().String("Type", "", "repository type (only YUM is supported)")
	cmd.Flags().String("SslCaCert", "", "SSL CA cert description, or an     empty string")
	cmd.Flags().String("SslCliCert", "", "SSL Client cert description, or     an empty string")
	cmd.Flags().String("SslCliKey", "", "SSL Client key description, or an     empty string")
	cmd.Flags().String("HasSignedMetadata", "", "true if the repository     has signed metadata, false otherwise")

	return cmd
}

func createRepo(globalFlags *types.GlobalFlags, flags *createRepoFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.Label, flags.Type, flags.Url, flags.Type, flags.SslCaCert, flags.SslCliCert, flags.SslCliKey, flags.Type, flags.SslCaCert, flags.SslCliCert, flags.SslCliKey, flags.HasSignedMetadata)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
