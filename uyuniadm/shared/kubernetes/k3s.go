package kubernetes

import (
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/templates"
)

const k3sTraefikConfigPath = "/var/lib/rancher/k3s/server/manifests/k3s-traefik-config.yaml"

func InstallK3sTraefikConfig() {
	log.Info().Msg("Installing K3s Traefik configuration")

	data := templates.K3sTraefikConfigTemplateData{
		TcpPorts: utils.TCP_PORTS,
		UdpPorts: utils.UDP_PORTS,
	}
	if err := utils.WriteTemplateToFile(data, k3sTraefikConfigPath, 0600, false); err != nil {
		log.Fatal().Err(err).Msgf("Failed to write K3s Traefik configuration")
	}

	// Wait for traefik to be back
	log.Info().Msg("Waiting for Traefik to be reloaded")
	for i := 0; i < 60; i++ {
		out, err := exec.Command("kubectl", "get", "job", "-A",
			"-o", "jsonpath={.status.completionTime}", "helm-install-traefik").Output()
		if err == nil {
			completionTime, err := time.Parse(time.RFC3339, string(out))
			if err == nil && time.Since(completionTime).Seconds() < 60 {
				break
			}
		}
	}
}

func UninstallK3sTraefikConfig(dryRun bool) {
	uninstallFile(k3sTraefikConfigPath, dryRun)
}
