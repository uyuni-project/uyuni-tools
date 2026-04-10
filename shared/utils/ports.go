// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewPortMap is a constructor for PortMap type.
func NewPortMap(port int) types.PortMap {
	return types.PortMap{
		Exposed: port,
		Port:    port,
	}
}

// WebPorts is the list of ports for the server web service.
var WebPorts = []types.PortMap{
	NewPortMap(80),
	NewPortMap(443),
}

// DBExporterPorts is the list of ports for the db exporter service.
var DBExporterPorts = []types.PortMap{
	NewPortMap(9187),
}

// ReportDBPorts is the list of ports for the server report db service.
var ReportDBPorts = []types.PortMap{
	NewPortMap(5432),
}

// DBPorts is the list of ports for the server internal db service.
var DBPorts = []types.PortMap{
	NewPortMap(5432),
}

// SaltPorts is the list of ports for the server salt service.
var SaltPorts = []types.PortMap{
	NewPortMap(4505),
	NewPortMap(4506),
}

// TaskoPorts is the list of ports for the server taskomatic service.
var TaskoPorts = []types.PortMap{
	NewPortMap(5556),
	NewPortMap(9800),
	NewPortMap(8001),
}

// TomcatPorts is the list of ports for the server tomcat service.
var TomcatPorts = []types.PortMap{
	NewPortMap(5557),
	NewPortMap(8003),
}

// SearchPorts is the list of ports for the server search service.
var SearchPorts = []types.PortMap{
	NewPortMap(8002),
}

// TftpPorts is the list of ports for the server tftp service.
var TftpPorts = []types.PortMap{
	{
		Exposed:  69,
		Port:     69,
		Protocol: "udp",
	},
}

const debugPortsStart = 8001
const debugPortsEnd = 8003

// GetServerPorts returns all the server container ports.
//
// if debug is set to true, the debug ports are added to the list.
func GetServerPorts(debug bool) []types.PortMap {
	ports := []types.PortMap{}
	ports = appendPorts(ports, debug, WebPorts...)
	ports = appendPorts(ports, debug, SaltPorts...)
	ports = appendPorts(ports, debug, TaskoPorts...)
	ports = appendPorts(ports, debug, TomcatPorts...)
	ports = appendPorts(ports, debug, SearchPorts...)
	ports = appendPorts(ports, debug, DBExporterPorts...)

	return ports
}

func appendPorts(ports []types.PortMap, debug bool, newPorts ...types.PortMap) []types.PortMap {
	for _, newPort := range newPorts {
		if debug || (newPort.Port < debugPortsStart || newPort.Port > debugPortsEnd) && !debug {
			ports = append(ports, newPort)
		}
	}
	return ports
}

// TCPPodmanPorts are the tcp ports required by the server on podman.
var TCPPodmanPorts = []types.PortMap{
	NewPortMap(9100),
}

// HubXmlrpcPorts are the tcp ports required by the Hub XMLRPC API service.
var HubXmlrpcPorts = []types.PortMap{
	NewPortMap(2830),
}

// GetProxyPorts returns all the proxy container ports.
func GetProxyPorts() []types.PortMap {
	ports := []types.PortMap{
		{
			Port:    22,
			Exposed: 8022,
		},
		NewPortMap(4505),
		NewPortMap(4506),
		NewPortMap(443),
		NewPortMap(80),
		{
			Exposed:  69,
			Port:     69,
			Protocol: "udp",
		},
	}

	return ports
}
