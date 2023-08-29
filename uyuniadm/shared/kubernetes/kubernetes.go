package kubernetes

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/templates"
)

type ClusterInfos struct {
	KubeletVersion string
	Ingress        string
}

func (infos ClusterInfos) IsK3s() bool {
	return strings.Contains(infos.KubeletVersion, "k3s")
}

func (infos ClusterInfos) IsRke2() bool {
	return strings.Contains(infos.KubeletVersion, "rke2")
}

func CheckCluster() ClusterInfos {

	// Get the kubelet version
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to get node hostname: %s\n", err)
	}

	out, err := exec.Command("kubectl", "get", "node",
		"-o", "jsonpath={.status.nodeInfo.kubeletVersion}", hostname).Output()
	if err != nil {
		log.Fatalf("Failed to get kubelet version for node %s: %s\n", hostname, err)
	}

	var infos ClusterInfos
	infos.KubeletVersion = string(out)
	infos.Ingress = guessIngress()

	return infos
}

func guessIngress() string {
	var ingress string

	// Check for a traefik resource
	err := exec.Command("kubectl", "explain", "ingressroutetcp").Run()
	if err == nil {
		ingress = "traefik"
	}

	// Look for a pod running the nginx-ingress-controller: there is no other common way to find out
	out, err := exec.Command("kubectl", "get", "pod", "-A",
		"-o", "jsonpath={range .items[*]}{.spec.containers[*].args[0]}{.spec.containers[*].command}{end}").Output()
	if err != nil {
		log.Fatalf("Failed to get get pod commands to look for nginx controller: %s", err)
	}

	const nginxController = "/nginx-ingress-controller"
	if strings.Contains(string(out), nginxController) {
		ingress = "nginx"
	}

	return ingress
}

const k3sTraefikConfigPath = "/var/lib/rancher/k3s/server/manifests/k3s-traefik-config.yaml"

func InstallK3sTraefikConfig() {
	log.Println("Installing K3s Traefik configuration")

	data := templates.K3sTraefikConfigTemplateData{
		TcpPorts: utils.TCP_PORTS,
		UdpPorts: utils.UDP_PORTS,
	}
	if err := utils.WriteTemplateToFile(data, k3sTraefikConfigPath, 0600, false); err != nil {
		log.Fatalf("Failed to write K3s Traefik configuration: %s\n", err)
	}

	// Wait for traefik to be back
	log.Println("Waiting for Traefik to be reloaded")
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

func uninstallFile(path string, dryRun bool) {
	if utils.FileExists(path) {
		if dryRun {
			log.Printf("Would remove file %s\n", path)
		} else {
			log.Printf("Removing file %s\n", path)
			if err := os.Remove(path); err != nil {
				log.Printf("Failed to remove file %s: %s\n", path, err)
			}
		}
	}
}
