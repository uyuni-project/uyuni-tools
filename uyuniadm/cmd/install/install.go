package install

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/install/kubernetes"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/install/podman"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	installCmd := &cobra.Command{
		Use:   "install [fqdn]",
		Short: "install a new server from scratch",
		Long: `Install a new server from scratch

The install command assumes the following:
  * podman or kubectl is installed locally
  * if kubectl is installed, a working kubeconfig should be set to connect to the cluster to deploy to

When installing on kubernetes, the helm values file will be overridden with the values from the uyuniadm parameters or configuration.

NOTE: for now installing on a remote cluster or podman is not supported!
`,
		Args: cobra.ExactArgs(1),
	}

	installCmd.AddCommand(podman.NewCommand(globalFlags))
	installCmd.AddCommand(kubernetes.NewCommand(globalFlags))

	return installCmd
}
