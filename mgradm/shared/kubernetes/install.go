// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// HELM_APP_NAME is the Helm application name.
const HELM_APP_NAME = "uyuni"

// Deploy execute a deploy of a given image and helm to a cluster.
func Deploy(cnx *shared.Connection, imageFlags *types.ImageFlags,
	helmFlags *cmd_utils.HelmFlags, sslFlags *cmd_utils.SslCertFlags, clusterInfos *kubernetes.ClusterInfos,
	fqdn string, debug bool, helmArgs ...string) error {
	// If installing on k3s, install the traefik helm config in manifests
	isK3s := clusterInfos.IsK3s()
	IsRke2 := clusterInfos.IsRke2()
	if isK3s {
		InstallK3sTraefikConfig(debug)
	} else if IsRke2 {
		kubernetes.InstallRke2NginxConfig(utils.TCP_PORTS, utils.UDP_PORTS, helmFlags.Uyuni.Namespace)
	}

	serverImage, err := utils.ComputeImage(imageFlags.Name, imageFlags.Tag)
	if err != nil {
		return fmt.Errorf("failed to compute image URL")
	}

	// Install the uyuni server helm chart
	err = UyuniUpgrade(serverImage, imageFlags.PullPolicy, helmFlags, clusterInfos.GetKubeconfig(), fqdn, clusterInfos.Ingress, helmArgs...)
	if err != nil {
		return fmt.Errorf("cannot upgrade: %s", err)
	}

	// Wait for the pod to be started
	err = kubernetes.WaitForDeployment(helmFlags.Uyuni.Namespace, HELM_APP_NAME, "uyuni")
	if err != nil {
		return fmt.Errorf("cannot deploy: %s", err)
	}
	return cnx.WaitForServer()
}

// DeployCertificate executre a deploy a new certificate given an helm.
func DeployCertificate(helmFlags *cmd_utils.HelmFlags, sslFlags *cmd_utils.SslCertFlags, rootCa string,
	ca *ssl.SslPair, kubeconfig string, fqdn string, imagePullPolicy string) ([]string, error) {
	helmArgs := []string{}
	if sslFlags.UseExisting() {
		DeployExistingCertificate(helmFlags, sslFlags, kubeconfig)
	} else {
		// Install cert-manager and a self-signed issuer ready for use
		issuerArgs, err := installSslIssuers(helmFlags, sslFlags, rootCa, ca, kubeconfig, fqdn, imagePullPolicy)
		if err != nil {
			return []string{}, fmt.Errorf("cannot install cert-manager and self-sign issuer: %s", err)
		}
		helmArgs = append(helmArgs, issuerArgs...)

		// Extract the CA cert into uyuni-ca config map as the container shouldn't have the CA secret
		extractCaCertToConfig()
	}

	return helmArgs, nil
}

// DeployExistingCertificate execute a deploy of an existing certificate.
func DeployExistingCertificate(helmFlags *cmd_utils.HelmFlags, sslFlags *cmd_utils.SslCertFlags, kubeconfig string) {
	// Deploy the SSL Certificate secret and CA configmap
	serverCrt, rootCaCrt := ssl.OrderCas(&sslFlags.Ca, &sslFlags.Server)
	serverKey := utils.ReadFile(sslFlags.Server.Key)
	installTlsSecret(helmFlags.Uyuni.Namespace, serverCrt, serverKey, rootCaCrt)

	// Extract the CA cert into uyuni-ca config map as the container shouldn't have the CA secret
	extractCaCertToConfig()
}

// UyuniUpgrade runs an helm upgrade using images and helm configuration as parameters.
func UyuniUpgrade(serverImage string, pullPolicy string, helmFlags *cmd_utils.HelmFlags, kubeconfig string,
	fqdn string, ingress string, helmArgs ...string) error {
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
	helmParams = append(helmParams,
		"--set", "images.server="+serverImage,
		"--set", "pullPolicy="+kubernetes.GetPullPolicy(pullPolicy),
		"--set", "fqdn="+fqdn)

	helmParams = append(helmParams, helmArgs...)

	namespace := helmFlags.Uyuni.Namespace
	chart := helmFlags.Uyuni.Chart
	version := helmFlags.Uyuni.Version
	return kubernetes.HelmUpgrade(kubeconfig, namespace, true, "", HELM_APP_NAME, chart, version, helmParams...)
}
