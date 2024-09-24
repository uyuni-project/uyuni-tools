// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
)

// deployNodeConfig deploy configuration files on the node.
func deployNodeConfig(
	namespace string,
	clusterInfos *kubernetes.ClusterInfos,
	needsHub bool,
	debug bool,
) error {
	// If installing on k3s, install the traefik helm config in manifests
	isK3s := clusterInfos.IsK3s()
	IsRke2 := clusterInfos.IsRke2()
	ports := getPortList(needsHub, debug)
	if isK3s {
		return kubernetes.InstallK3sTraefikConfig(ports)
	} else if IsRke2 {
		return kubernetes.InstallRke2NginxConfig(ports, namespace)
	}
	return nil
}
