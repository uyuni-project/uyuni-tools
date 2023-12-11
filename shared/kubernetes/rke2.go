// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const rke2NginxConfigPath = "/var/lib/rancher/rke2/server/manifests/rke2-ingress-nginx-config.yaml"

func InstallRke2NginxConfig(tcpPorts []types.PortMap, udpPorts []types.PortMap, namespace string) {
	log.Info().Msg("Installing RKE2 Nginx configuration")

	data := Rke2NginxConfigTemplateData{
		Namespace: namespace,
		TcpPorts:  tcpPorts,
		UdpPorts:  udpPorts,
	}
	if err := utils.WriteTemplateToFile(data, rke2NginxConfigPath, 0600, false); err != nil {
		log.Fatal().Err(err).Msgf("Failed to write Rke2 nginx configuration")
	}

	// Wait for the nginx controller to be back
	log.Info().Msg("Waiting for Nginx controller to be reloaded")
	for i := 0; i < 60; i++ {
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "daemonset", "-A",
			"-o", "jsonpath={.status.numberReady}", "rke2-ingress-nginx-controller")
		if err == nil {
			if count, err := strconv.Atoi(string(out)); err == nil && count > 0 {
				break
			}
		}
	}
}

func UninstallRke2NginxConfig(dryRun bool) {
	utils.UninstallFile(rke2NginxConfigPath, dryRun)
}
