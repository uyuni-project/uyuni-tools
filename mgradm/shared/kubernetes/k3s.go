// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const k3sTraefikConfigPath = "/var/lib/rancher/k3s/server/manifests/k3s-traefik-config.yaml"

func InstallK3sTraefikConfig(debug bool) {
	tcpPorts := []types.PortMap{}
	tcpPorts = append(tcpPorts, utils.TCP_PORTS...)
	if debug {
		tcpPorts = append(tcpPorts, utils.DEBUG_PORTS...)
	}

	kubernetes.InstallK3sTraefikConfig(tcpPorts, utils.UDP_PORTS)
}
