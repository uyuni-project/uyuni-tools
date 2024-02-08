// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"path"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
)

const helmAppName = "uyuni-proxy"

func Deploy(installFlags *utils.ProxyInstallFlags, helmFlags *HelmFlags, configDir string,
	kubeconfig string, helmArgs ...string,
) error {
	log.Info().Msg("Installing Uyuni")

	helmParams := []string{}

	// Pass the user-provided values file
	extraValues := helmFlags.Proxy.Values
	if extraValues != "" {
		helmParams = append(helmParams, "-f", extraValues)
	}

	helmParams = append(helmParams,
		"-f", path.Join(configDir, "httpd.yaml"),
		"-f", path.Join(configDir, "ssh.yaml"),
		"-f", path.Join(configDir, "config.yaml"),
		"--set", "images.proxy-httpd="+installFlags.GetContainerImage("httpd"),
		"--set", "images.proxy-salt-broker="+installFlags.GetContainerImage("salt-broker"),
		"--set", "images.proxy-squid="+installFlags.GetContainerImage("squid"),
		"--set", "images.proxy-ssh="+installFlags.GetContainerImage("ssh"),
		"--set", "images.proxy-tftpd="+installFlags.GetContainerImage("tftpd"),
		"--set", "repository="+installFlags.ImagesLocation,
		"--set", "version="+installFlags.Tag,
		"--set", "pullPolicy="+kubernetes.GetPullPolicy(installFlags.PullPolicy))

	helmParams = append(helmParams, helmArgs...)

	// Install the helm chart
	kubernetes.HelmUpgrade(kubeconfig, helmFlags.Proxy.Namespace, true, "", helmAppName, helmFlags.Proxy.Chart,
		helmFlags.Proxy.Version, helmParams...)

	// Wait for the pod to be started
	return kubernetes.WaitForDeployment(helmFlags.Proxy.Namespace, helmAppName, "uyuni-proxy")
}
