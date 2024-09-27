// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import "github.com/uyuni-project/uyuni-tools/shared/types"

// ServerTCPServiceName is the name of the server TCP service.
const ServerTCPServiceName = "uyuni-tcp"

// ServerUDPServiceName is the name of the server UDP service.
const ServerUDPServiceName = "uyuni-udp"

// ProxyTCPServiceName is the name of the proxy TCP service.
const ProxyTCPServiceName = "uyuni-proxy-tcp"

// ProxyUDPServiceName is the name of the proxy UDP service.
const ProxyUDPServiceName = "uyuni-proxy-udp"

// NewPortMap is a constructor for PortMap type.
func NewPortMap(service string, name string, exposed int, port int) types.PortMap {
	return types.PortMap{
		Service: service,
		Name:    name,
		Exposed: exposed,
		Port:    port,
	}
}

// WebPorts is the list of ports for the server web service.
var WebPorts = []types.PortMap{
	NewPortMap(ServerTCPServiceName, "http", 80, 80),
}

// PgsqlPorts is the list of ports for the server report db service.
var PgsqlPorts = []types.PortMap{
	NewPortMap(ServerTCPServiceName, "pgsql", 5432, 5432),
	NewPortMap(ServerTCPServiceName, "exporter", 9187, 9187),
}

// SaltPorts is the list of ports for the server salt service.
var SaltPorts = []types.PortMap{
	NewPortMap(ServerTCPServiceName, "publish", 4505, 4505),
	NewPortMap(ServerTCPServiceName, "request", 4506, 4506),
}

// CobblerPorts is the list of ports for the server cobbler service.
var CobblerPorts = []types.PortMap{
	NewPortMap(ServerTCPServiceName, "cobbler", 25151, 25151),
}

// TaskoPorts is the list of ports for the server taskomatic service.
var TaskoPorts = []types.PortMap{
	NewPortMap(ServerTCPServiceName, "jmx", 5556, 5556),
	NewPortMap(ServerTCPServiceName, "mtrx", 9800, 9800),
	NewPortMap(ServerTCPServiceName, "debug", 8001, 8001),
}

// TomcatPorts is the list of ports for the server tomcat service.
var TomcatPorts = []types.PortMap{
	NewPortMap(ServerTCPServiceName, "jmx", 5557, 5557),
	NewPortMap(ServerTCPServiceName, "debug", 8003, 8003),
}

// SearchPorts is the list of ports for the server search service.
var SearchPorts = []types.PortMap{
	NewPortMap(ServerTCPServiceName, "debug", 8002, 8002),
}

// TftpPorts is the list of ports for the server tftp service.
var TftpPorts = []types.PortMap{
	{
		Service:  ServerUDPServiceName,
		Name:     "tftp",
		Exposed:  69,
		Port:     69,
		Protocol: "udp",
	},
}

// GetServerPorts returns all the server container ports.
//
// if debug is set to true, the debug ports are added to the list.
func GetServerPorts(debug bool) []types.PortMap {
	ports := []types.PortMap{}
	ports = appendPorts(ports, debug, WebPorts...)
	ports = appendPorts(ports, debug, PgsqlPorts...)
	ports = appendPorts(ports, debug, SaltPorts...)
	ports = appendPorts(ports, debug, CobblerPorts...)
	ports = appendPorts(ports, debug, TaskoPorts...)
	ports = appendPorts(ports, debug, TomcatPorts...)
	ports = appendPorts(ports, debug, SearchPorts...)
	ports = appendPorts(ports, debug, TftpPorts...)

	return ports
}

func appendPorts(ports []types.PortMap, debug bool, newPorts ...types.PortMap) []types.PortMap {
	for _, newPort := range newPorts {
		if debug || newPort.Name != "debug" && !debug {
			ports = append(ports, newPort)
		}
	}
	return ports
}

// TCPPodmanPorts are the tcp ports required by the server on podman.
var TCPPodmanPorts = []types.PortMap{
	// TODO: Replace Node exporter with cAdvisor
	NewPortMap("tomcat", "node-exporter", 9100, 9100),
}

// HubXmlrpcPorts are the tcp ports required by the Hub XMLRPC API service.
var HubXmlrpcPorts = []types.PortMap{
	NewPortMap(ServerTCPServiceName, "xmlrpc", 2830, 2830),
}

// ProxyTCPPorts are the tcp ports required by the proxy.
var ProxyTCPPorts = []types.PortMap{
	NewPortMap(ProxyTCPServiceName, "ssh", 8022, 22),
	NewPortMap(ProxyTCPServiceName, "publish", 4505, 4505),
	NewPortMap(ProxyTCPServiceName, "request", 4506, 4506),
}

// ProxyPodmanPorts are the http/s ports required by the proxy.
var ProxyPodmanPorts = []types.PortMap{
	NewPortMap(ProxyTCPServiceName, "https", 443, 443),
	NewPortMap(ProxyTCPServiceName, "http", 80, 80),
}

// GetProxyPorts returns all the proxy container ports.
func GetProxyPorts() []types.PortMap {
	ports := []types.PortMap{}
	ports = appendPorts(ports, false, ProxyTCPPorts...)
	ports = appendPorts(ports, false, types.PortMap{
		Service:  ProxyUDPServiceName,
		Name:     "tftp",
		Exposed:  69,
		Port:     69,
		Protocol: "udp",
	})

	return ports
}
