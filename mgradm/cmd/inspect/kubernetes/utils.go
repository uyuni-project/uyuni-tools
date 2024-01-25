// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

//import (
//	"os/exec"
//
//	"github.com/rs/zerolog"
//	"github.com/rs/zerolog/log"
//	"github.com/spf13/cobra"
//	install_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/shared"
//	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
//	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
//	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
//	"github.com/uyuni-project/uyuni-tools/shared"
//	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
//	"github.com/uyuni-project/uyuni-tools/shared/types"
//)
//
//func installForKubernetes(globalFlags *types.GlobalFlags,
//	flags *kubernetesInstallFlags,
//	cmd *cobra.Command,
//	args []string,
//) error {
//	for _, binary := range []string{"kubectl", "helm"} {
//		if _, err := exec.LookPath(binary); err != nil {
//			log.Fatal().Err(err).Msgf("install %s before running this command", binary)
//		}
//	}
//
//	flags.CheckParameters(cmd, "kubectl")
//	cnx := shared.NewConnection("kubectl", "", shared_kubernetes.ServerFilter)
//
//	fqdn := args[0]
//
//	helmArgs := []string{"--set", "timezone=" + flags.TZ}
//	if flags.MirrorPath != "" {
//		// TODO Handle claims for multi-node clusters
//		helmArgs = append(helmArgs, "--set", "mirror.hostPath="+flags.MirrorPath)
//	}
//	if flags.Debug.Java {
//		helmArgs = append(helmArgs, "--set", "exposeJavaDebug=true")
//	}
//
//	// Check the kubernetes cluster setup
//	clusterInfos := shared_kubernetes.CheckCluster()
//
//	// Deploy the SSL CA or server certificate
//	ca := ssl.SslPair{}
//	sslArgs := kubernetes.DeployCertificate(&flags.Helm, &flags.Ssl, "", &ca, clusterInfos.GetKubeconfig(), fqdn,
//		flags.Image.PullPolicy)
//	helmArgs = append(helmArgs, sslArgs...)
//
//	// Deploy Uyuni and wait for it to be up
//	kubernetes.Deploy(cnx, &flags.Image, &flags.Helm, &flags.Ssl, &clusterInfos, fqdn, flags.Debug.Java, helmArgs...)
//
//	// Create setup script + env variables and copy it to the container
//	envs := map[string]string{
//		"NO_SSL": "Y",
//	}
//
//	install_shared.RunSetup(cnx, &flags.InstallFlags, args[0], envs)
//
//	// The CA needs to be added to the database for Kickstart use.
//	err := adm_utils.ExecCommand(zerolog.DebugLevel, cnx,
//		"/usr/bin/rhn-ssl-dbstore", "--ca-cert=/etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT")
//	if err != nil {
//		log.Fatal().Err(err).Msg("Error storing the SSL CA certificate in database")
//	}
//	return nil
//}
