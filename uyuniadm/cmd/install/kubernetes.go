package install

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/kubernetes"
)

func installForKubernetes(globalFlags *types.GlobalFlags, flags *InstallFlags, cmd *cobra.Command, args []string) {
	fqdn := args[0]

	helmArgs := []string{"--set", "timezone=" + flags.TZ}
	if flags.MirrorPath != "" {
		// TODO Handle claims for multi-node clusters
		helmArgs = append(helmArgs, "--set", "mirror.hostPath="+flags.MirrorPath)
	}

	// Check the kubernetes cluster setup
	clusterInfos := kubernetes.CheckCluster()

	// Deploy the SSL CA or server certificate
	ca := kubernetes.TlsCert{}
	sslArgs := kubernetes.DeployCertificate(globalFlags, &flags.Helm, &flags.Cert, &ca, clusterInfos.GetKubeconfig(), fqdn)
	helmArgs = append(helmArgs, sslArgs...)

	// Deploy Uyuni and wait for it to be up
	kubernetes.Deploy(globalFlags, &flags.Image, &flags.Helm, &flags.Cert, &clusterInfos, fqdn, helmArgs...)

	// Create setup script + env variables and copy it to the container
	envs := map[string]string{
		"NO_SSL": "Y",
	}

	runSetup(globalFlags, flags, args[0], envs)
}
