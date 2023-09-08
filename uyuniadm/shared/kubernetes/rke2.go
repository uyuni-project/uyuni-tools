package kubernetes

import (
	"log"
	"os/exec"
	"strconv"

	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/templates"
)

const rke2NginxConfigPath = "/var/lib/rancher/rke2/server/manifests/rke2-ingress-nginx-config.yaml"

func InstallRke2NginxConfig(namespace string) {
	log.Println("Installing RKE2 Nginx configuration")

	data := templates.Rke2NginxConfigTemplateData{
		Namespace: namespace,
		TcpPorts:  utils.TCP_PORTS,
		UdpPorts:  utils.UDP_PORTS,
	}
	if err := utils.WriteTemplateToFile(data, rke2NginxConfigPath, 0600, false); err != nil {
		log.Fatalf("Failed to write Rke2 nginx configuration: %s\n", err)
	}

	// Wait for the nginx controller to be back
	log.Println("Waiting for Nginx controller to be reloaded")
	for i := 0; i < 60; i++ {
		out, err := exec.Command("kubectl", "get", "daemonset", "-A",
			"-o", "jsonpath={.status.numberReady}", "rke2-ingress-nginx-controller").Output()
		if err == nil {
			if count, err := strconv.Atoi(string(out)); err == nil && count > 0 {
				break
			}
		}
	}
}

func UninstallRke2NginxConfig(dryRun bool) {
	uninstallFile(rke2NginxConfigPath, dryRun)
}
