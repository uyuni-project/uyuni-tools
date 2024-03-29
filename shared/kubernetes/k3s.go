// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const k3sTraefikConfigPath = "/var/lib/rancher/k3s/server/manifests/k3s-traefik-config.yaml"

// InstallK3sTraefikConfig install K3s Traefik configuration.
func InstallK3sTraefikConfig(tcpPorts []types.PortMap, udpPorts []types.PortMap) {
	log.Info().Msg("Installing K3s Traefik configuration")

	data := K3sTraefikConfigTemplateData{
		TcpPorts: tcpPorts,
		UdpPorts: udpPorts,
	}
	if err := utils.WriteTemplateToFile(data, k3sTraefikConfigPath, 0600, false); err != nil {
		log.Fatal().Err(err).Msgf("Failed to write K3s Traefik configuration")
	}

	// Wait for traefik to be back
	log.Info().Msg("Waiting for Traefik to be reloaded")
	for i := 0; i < 60; i++ {
		out, err := utils.RunCmdOutput(zerolog.TraceLevel, "kubectl", "get", "job", "-A",
			"-o", "jsonpath={.status.completionTime}", "helm-install-traefik")
		if err == nil {
			completionTime, err := time.Parse(time.RFC3339, string(out))
			if err == nil && time.Since(completionTime).Seconds() < 60 {
				break
			}
		}
	}
}

// UninstallK3sTraefikConfig uninstall K3s Traefik configuration.
func UninstallK3sTraefikConfig(dryRun bool) {
	utils.UninstallFile(k3sTraefikConfigPath, dryRun)
}
