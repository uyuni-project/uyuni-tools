// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// ClusterInfos represent cluster information.
type ClusterInfos struct {
	KubeletVersion string
	Ingress        string
}

// IsK3s is true if it's a K3s Cluster.
func (infos ClusterInfos) IsK3s() bool {
	return strings.Contains(infos.KubeletVersion, "k3s")
}

// IsRKE2 is true if it's a RKE2 Cluster.
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

// CheckCluster return cluster information.
func CheckCluster() ClusterInfos {
	// Get the kubelet version
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to get node hostname")
	}

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "node",
		"-o", "jsonpath={.status.nodeInfo.kubeletVersion}", hostname)
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
	err := utils.RunCmd("kubectl", "explain", "ingressroutetcp")
	if err == nil {
		ingress = "traefik"
	}

	// Look for a pod running the nginx-ingress-controller: there is no other common way to find out
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "pod", "-A",
		"-o", "jsonpath={range .items[*]}{.spec.containers[*].args[0]}{.spec.containers[*].command}{end}")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to get pod commands to look for nginx controller")
	}

	const nginxController = "/nginx-ingress-controller"
	if strings.Contains(string(out), nginxController) {
		ingress = "nginx"
	}

	return ingress
}

// Restart restarts the pod.
func Restart(ServerFilter string) error {
	if err := Stop(ServerFilter); err != nil {
		return fmt.Errorf("cannot stop %s: %s", ServerFilter, err)
	}
	return Start(ServerFilter)
}

// Start starts the pod.
func Start(ServerFilter string) error {
	// if something is running, we don't need to set replicas to 1
	if _, err := GetNode("uyuni"); err != nil {
		return ReplicasTo(ServerFilter, 1)
	}
	log.Debug().Msgf("Already running")
	return nil
}

// Stop stop the pod.
func Stop(ServerFilter string) error {
	return ReplicasTo(ServerFilter, 0)
}
