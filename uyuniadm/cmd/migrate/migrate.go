package migrate

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/migrate/kubernetes"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/migrate/podman"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	migrateCmd := &cobra.Command{
		Use:   "migrate [source server FQDN]",
		Short: "migrate a remote server to containers",
		Long: `Migrate a remote server to containers

This migration command assumes a few things:
  * the SSH configuration for the source server is complete, including user and
    all needed options to connect to the machine,
  * an SSH agent is started and the key to use to connect to the server is added to it,
  * podman or kubectl is installed locally
  * if kubectl is installed, a working kubeconfig should be set to connect to the cluster to deploy to

NOTE: for now installing on a remote cluster or podman is not supported yet!
`,
	}

	migrateCmd.AddCommand(podman.NewCommand(globalFlags))
	migrateCmd.AddCommand(kubernetes.NewCommand(globalFlags))

	return migrateCmd
}
