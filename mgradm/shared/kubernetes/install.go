// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const HELM_APP_NAME = "uyuni"

func Deploy(cnx *utils.Connection, imageFlags *types.ImageFlags,
	helmFlags *cmd_utils.HelmFlags, sslFlags *cmd_utils.SslCertFlags, clusterInfos *kubernetes.ClusterInfos,
	fqdn string, debug bool, helmArgs ...string) {

	// If installing on k3s, install the traefik helm config in manifests
	isK3s := clusterInfos.IsK3s()
	IsRke2 := clusterInfos.IsRke2()
	if isK3s {
		InstallK3sTraefikConfig(debug)
	} else if IsRke2 {
		kubernetes.InstallRke2NginxConfig(utils.TCP_PORTS, utils.UDP_PORTS, helmFlags.Uyuni.Namespace)
	}

	// Install the uyuni server helm chart
	UyuniUpgrade(imageFlags, helmFlags, clusterInfos.GetKubeconfig(), fqdn, clusterInfos.Ingress, helmArgs...)

	// Wait for the pod to be started
	kubernetes.WaitForDeployment(helmFlags.Uyuni.Namespace, HELM_APP_NAME, "uyuni")
	cnx.WaitForServer()
}

func DeployCertificate(helmFlags *cmd_utils.HelmFlags, sslFlags *cmd_utils.SslCertFlags, rootCa string,
	ca *ssl.SslPair, kubeconfig string, fqdn string, imagePullPolicy string) []string {

	helmArgs := []string{}
	if sslFlags.UseExisting() {
		DeployExistingCertificate(helmFlags, sslFlags, kubeconfig)
	} else {
		// Install cert-manager and a self-signed issuer ready for use
		issuerArgs := installSslIssuers(helmFlags, sslFlags, rootCa, ca, kubeconfig, fqdn, imagePullPolicy)
		helmArgs = append(helmArgs, issuerArgs...)

		// Extract the CA cert into uyuni-ca config map as the container shouldn't have the CA secret
		extractCaCertToConfig()
	}

	return helmArgs
}

func DeployExistingCertificate(helmFlags *cmd_utils.HelmFlags, sslFlags *cmd_utils.SslCertFlags, kubeconfig string) {

	// Deploy the SSL Certificate secret and CA configmap
	serverCrt, rootCaCrt := ssl.OrderCas(&sslFlags.Ca, &sslFlags.Server)
	serverKey := utils.ReadFile(sslFlags.Server.Key)
	installTlsSecret(helmFlags.Uyuni.Namespace, serverCrt, serverKey, rootCaCrt)

	// Extract the CA cert into uyuni-ca config map as the container shouldn't have the CA secret
	extractCaCertToConfig()
}

func UyuniUpgrade(imageFlags *types.ImageFlags, helmFlags *cmd_utils.HelmFlags, kubeconfig string,
	fqdn string, ingress string, helmArgs ...string) {

	log.Info().Msg("Installing Uyuni")

	// The guessed ingress is passed before the user's value to let the user override it in case we got it wrong.
	helmParams := []string{
		"--set", "ingress=" + ingress,
	}

	extraValues := helmFlags.Uyuni.Values
	if extraValues != "" {
		helmParams = append(helmParams, "-f", extraValues)
	}

	// The values computed from the command line need to be last to override what could be in the extras
	serverImage, err := utils.ComputeImage(imageFlags.Name, imageFlags.Tag)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to compute image URL")
	}
	helmParams = append(helmParams,
		"--set", serverImage,
		"--set", "pullPolicy="+kubernetes.GetPullPolicy(imageFlags.PullPolicy),
		"--set", "fqdn="+fqdn)

	helmParams = append(helmParams, helmArgs...)

	namespace := helmFlags.Uyuni.Namespace
	chart := helmFlags.Uyuni.Chart
	version := helmFlags.Uyuni.Version
	kubernetes.HelmUpgrade(kubeconfig, namespace, true, "", HELM_APP_NAME, chart, version, helmParams...)
}
