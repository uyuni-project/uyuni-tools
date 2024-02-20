// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// InstallK3sTraefikConfig installs the K3s Traefik configuration.
func InstallK3sTraefikConfig(debug bool) {
	tcpPorts := []types.PortMap{}
	tcpPorts = append(tcpPorts, utils.TCP_PORTS...)
	if debug {
		tcpPorts = append(tcpPorts, utils.DEBUG_PORTS...)
	}

	kubernetes.InstallK3sTraefikConfig(tcpPorts, utils.UDP_PORTS)
}
