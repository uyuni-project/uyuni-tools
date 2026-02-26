// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// ClusterInfos represent cluster information.
type ClusterInfos struct {
	KubeletVersion string
}

// IsK3s is true if it's a K3s Cluster.
func (infos ClusterInfos) IsK3s() bool {
	return strings.Contains(infos.KubeletVersion, "k3s")
}

// IsRke2 is true if it's a RKE2 Cluster.
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
func CheckCluster() (*ClusterInfos, error) {
	// Get the kubelet version
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "node",
		"-o", "jsonpath={.items[0].status.nodeInfo.kubeletVersion}")
	if err != nil {
		return nil, utils.Errorf(err, L("failed to get kubelet version"))
	}

	var infos ClusterInfos
	infos.KubeletVersion = string(out)
	if err != nil {
		return nil, err
	}

	return &infos, nil
}

// Restart restarts the pod.
func Restart(namespace string, app string) error {
	if err := Stop(namespace, app); err != nil {
		return utils.Errorf(err, L("cannot stop %s"), app)
	}
	return Start(namespace, app)
}

// Start starts the pod.
func Start(namespace string, app string) error {
	// if something is running, we don't need to set replicas to 1
	if _, err := GetNode(namespace, "-l"+AppLabel+"="+app); err != nil {
		return ReplicasTo(namespace, app, 1)
	}
	log.Debug().Msgf("Already running")
	return nil
}

// Stop stop the pod.
func Stop(namespace string, app string) error {
	return ReplicasTo(namespace, app, 0)
}
