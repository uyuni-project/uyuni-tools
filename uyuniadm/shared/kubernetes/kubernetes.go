package kubernetes

import (
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
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

// GetKubeconfig returns the path to the default kubeconfig file or "" if none.
func (infos ClusterInfos) GetKubeconfig() string {
	var kubeconfig string
	if infos.IsK3s() {
		// If the user didn't provide a KUBECONFIG value or file, use the k3s default
		kubeconfigPath := os.ExpandEnv("${HOME}/.kube/config")
		if os.Getenv("KUBECONFIG") == "" || !utils.FileExists(kubeconfigPath) {
			kubeconfig = "/etc/rancher/k3s/k3s.yaml"
		}
	}
	// Since even kubectl doesn't work without a trick on rke2, we assume the user has set kubeconfig
	return kubeconfig
}

func CheckCluster() ClusterInfos {

	// Get the kubelet version
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to get node hostname")
	}

	out, err := exec.Command("kubectl", "get", "node",
		"-o", "jsonpath={.status.nodeInfo.kubeletVersion}", hostname).Output()
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to get kubelet version for node %s", hostname)
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
		log.Fatal().Err(err).Msgf("Failed to get get pod commands to look for nginx controller")
	}

	const nginxController = "/nginx-ingress-controller"
	if strings.Contains(string(out), nginxController) {
		ingress = "nginx"
	}

	return ingress
}
