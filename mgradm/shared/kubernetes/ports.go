// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// getPortList returns compiled lists of tcp and udp ports..
func getPortList(hub bool, debug bool) []types.PortMap {
	ports := utils.GetServerPorts(debug)
	ports = append(ports, utils.ReportDBPorts...)

	return ports
}
