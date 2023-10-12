//go:build !nok8s

package kubernetes

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/install/shared"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/ssl"
	adm_utils "github.com/uyuni-project/uyuni-tools/uyuniadm/shared/utils"
)

func installForKubernetes(globalFlags *types.GlobalFlags, flags *kubernetesInstallFlags,
	cmd *cobra.Command, args []string) {
	cnx := utils.NewConnection("kubectl")

	fqdn := args[0]

	helmArgs := []string{"--set", "timezone=" + flags.TZ}
	if flags.MirrorPath != "" {
		// TODO Handle claims for multi-node clusters
		helmArgs = append(helmArgs, "--set", "mirror.hostPath="+flags.MirrorPath)
	}
	if flags.Debug.Java {
		helmArgs = append(helmArgs, "--set", "exposeJavaDebug=true")
	}

	// Check the kubernetes cluster setup
	clusterInfos := kubernetes.CheckCluster()

	// Deploy the SSL CA or server certificate
	ca := ssl.SslPair{}
	sslArgs := kubernetes.DeployCertificate(&flags.Helm, &flags.Ssl, "", &ca, clusterInfos.GetKubeconfig(), fqdn)
	helmArgs = append(helmArgs, sslArgs...)

	// Deploy Uyuni and wait for it to be up
	kubernetes.Deploy(cnx, &flags.Image, &flags.Helm, &flags.Ssl, &clusterInfos, fqdn, flags.Debug.Java, helmArgs...)

	// Create setup script + env variables and copy it to the container
	envs := map[string]string{
		"NO_SSL": "Y",
	}

	shared.RunSetup(cnx, &flags.InstallFlags, args[0], envs)

	// The CA needs to be added to the database for Kickstart use.
	err := adm_utils.ExecCommand(zerolog.DebugLevel, cnx,
		"/usr/bin/rhn-ssl-dbstore", "--ca-cert=/etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT")
	if err != nil {
		log.Fatal().Err(err).Msg("Error storing the SSL CA certificate in database")
	}
}
