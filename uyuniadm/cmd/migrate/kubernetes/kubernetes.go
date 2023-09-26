//go:build !nok8s

package kubernetes

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/migrate/shared"
	cmd_utils "github.com/uyuni-project/uyuni-tools/uyuniadm/shared/utils"
)

type kubernetesMigrateFlags struct {
	shared.MigrateFlags `mapstructure:",squash"`
	Helm                cmd_utils.HelmFlags
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	migrateCmd := &cobra.Command{
		Use:   "kubernetes [source server FQDN]",
		Short: "migrate a remote server to containers running on a kubernetes cluster",
		Long: `Migrate a remote server to containers running on a kubernetes cluster

This migration command assumes a few things:
  * the SSH configuration for the source server is complete, including user and
    all needed options to connect to the machine,
  * an SSH agent is started and the key to use to connect to the server is added to it,
  * kubectl is installed locally
  * A working kubeconfig should be set to connect to the cluster to deploy to

NOTE: for now installing on a remote cluster is not supported yet!
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viper := utils.ReadConfig(globalFlags.ConfigPath, "admconfig", cmd)
			var flags kubernetesMigrateFlags
			if err := viper.Unmarshal(&flags); err != nil {
				log.Fatal().Err(err).Msg("Failed to Unmarshal configuration")
			}

			migrateToKubernetes(globalFlags, &flags, cmd, args)
		},
	}

	shared.AddMigrateFlags(migrateCmd)
	cmd_utils.AddHelmInstallFlag(migrateCmd)

	return migrateCmd
}
