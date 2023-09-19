package install

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
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

	// The CA needs to be added to the database for Kickstart use.
	err := utils.Exec(globalFlags, "kubectl", false, false, true, []string{},
		"/usr/bin/rhn-ssl-dbstore", "--ca-cert=/etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT")
	if err != nil {
		log.Fatal().Err(err).Msg("Error storing the SSL CA certificate in database")
	}
}
