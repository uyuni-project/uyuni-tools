// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const rke2NginxConfigPath = "/var/lib/rancher/rke2/server/manifests/uyuni-ingress-nginx-config.yaml"

// InstallRke2NgixConfig install Rke2 Nginx configuration.
func InstallRke2NginxConfig(ports []types.PortMap, namespace string) error {
	log.Info().Msg(L("Installing RKE2 Nginx configuration"))

	tcpPorts := []types.PortMap{}
	udpPorts := []types.PortMap{}
	for _, port := range ports {
		if port.Protocol == "udp" {
			udpPorts = append(udpPorts, port)
		} else {
			tcpPorts = append(tcpPorts, port)
		}
	}

	data := Rke2NginxConfigTemplateData{
		Namespace: namespace,
		TCPPorts:  tcpPorts,
		UDPPorts:  udpPorts,
	}
	if err := utils.WriteTemplateToFile(data, rke2NginxConfigPath, 0o600, true); err != nil {
		return utils.Errorf(err, L("Failed to write Rke2 nginx configuration"))
	}

	// Wait for the nginx controller to be back
	log.Info().Msg(L("Waiting for Nginx controller to be reloaded"))
	for i := 0; i < 60; i++ {
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "daemonset", "-A",
			"-o", "jsonpath={.status.numberReady}", "rke2-ingress-nginx-controller")
		if err == nil {
			if count, err := strconv.Atoi(string(out)); err == nil && count > 0 {
				break
			}
		}
	}
	return nil
}

// UninstallRke2NginxConfig uninstall Rke2 Nginx configuration.
func UninstallRke2NginxConfig(dryRun bool) {
	utils.UninstallFile(rke2NginxConfigPath, dryRun)
}
