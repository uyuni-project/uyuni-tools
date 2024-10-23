// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

const (
	// WebServiceName is the name of the server web service.
	WebServiceName = "web"
	// SaltServiceName is the name of the server salt service.
	SaltServiceName = "salt"
	// CobblerServiceName is the name of the server cobbler service.
	CobblerServiceName = "cobbler"
	// ReportdbServiceName is the name of the server report database service.
	ReportdbServiceName = "reportdb"
	// DbServiceName is the name of the server internal database service.
	DbServiceName = "db"
	// TaskoServiceName is the name of the server taskomatic service.
	TaskoServiceName = "taskomatic"
	// TftpServiceName is the name of the server tftp service.
	TftpServiceName = "tftp"
	// TomcatServiceName is the name of the server tomcat service.
	TomcatServiceName = "tomcat"
	// SearchServiceName is the name of the server search service.
	SearchServiceName = "search"

	// HubApiServiceName is the name of the server hub API service.
	HubApiServiceName = "hub-api"

	// ProxyTcpServiceName is the name of the proxy TCP service.
	ProxyTcpServiceName = "uyuni-proxy-tcp"

	// ProxyUdpServiceName is the name of the proxy UDP service.
	ProxyUdpServiceName = "uyuni-proxy-udp"
)

// NewPortMap is a constructor for PortMap type.
func NewPortMap(service string, name string, exposed int, port int) types.PortMap {
	return types.PortMap{
		Service: service,
		Name:    name,
		Exposed: exposed,
		Port:    port,
	}
}

// WEB_PORTS is the list of ports for the server web service.
var WEB_PORTS = []types.PortMap{
	NewPortMap(WebServiceName, "http", 80, 80),
}

// REPORTDB_PORTS is the list of ports for the server report db service.
var REPORTDB_PORTS = []types.PortMap{
	NewPortMap(ReportdbServiceName, "pgsql", 5432, 5432),
	NewPortMap(ReportdbServiceName, "exporter", 9187, 9187),
}

// DB_PORTS is the list of ports for the server internal db service.
var DB_PORTS = []types.PortMap{
	NewPortMap(DbServiceName, "pgsql", 5432, 5432),
	NewPortMap(DbServiceName, "exporter", 9187, 9187),
}

// SALT_PORTS is the list of ports for the server salt service.
var SALT_PORTS = []types.PortMap{
	NewPortMap(SaltServiceName, "publish", 4505, 4505),
	NewPortMap(SaltServiceName, "request", 4506, 4506),
}

// COBBLER_PORTS is the list of ports for the server cobbler service.
var COBBLER_PORTS = []types.PortMap{
	NewPortMap(CobblerServiceName, "cobbler", 25151, 25151),
}

// TASKO_PORTS is the list of ports for the server taskomatic service.
var TASKO_PORTS = []types.PortMap{
	NewPortMap(TaskoServiceName, "jmx", 5556, 5556),
	NewPortMap(TaskoServiceName, "mtrx", 9800, 9800),
	NewPortMap(TaskoServiceName, "debug", 8001, 8001),
}

// TOMCAT_PORTS is the list of ports for the server tomcat service.
var TOMCAT_PORTS = []types.PortMap{
	NewPortMap(TomcatServiceName, "jmx", 5557, 5557),
	NewPortMap(TomcatServiceName, "debug", 8003, 8003),
}

// SEARCH_PORTS is the list of ports for the server search service.
var SEARCH_PORTS = []types.PortMap{
	NewPortMap(SearchServiceName, "debug", 8002, 8002),
}

// TFTP_PORTS is the list of ports for the server tftp service.
var TFTP_PORTS = []types.PortMap{
	{
		Service:  TftpServiceName,
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
	ports = appendPorts(ports, debug, WEB_PORTS...)
	ports = appendPorts(ports, debug, REPORTDB_PORTS...)
	// DB_PORTS are not added here as they are not exposed and not belonging to the server soon.
	ports = appendPorts(ports, debug, SALT_PORTS...)
	ports = appendPorts(ports, debug, COBBLER_PORTS...)
	ports = appendPorts(ports, debug, TASKO_PORTS...)
	ports = appendPorts(ports, debug, TOMCAT_PORTS...)
	ports = appendPorts(ports, debug, SEARCH_PORTS...)
	ports = appendPorts(ports, debug, TFTP_PORTS...)

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

// TCP_PODMAN_PORTS are the tcp ports required by the server on podman.
var TCP_PODMAN_PORTS = []types.PortMap{
	// TODO: Replace Node exporter with cAdvisor
	NewPortMap("tomcat", "node-exporter", 9100, 9100),
}

// HUB_XMLRPC_PORTS are the tcp ports required by the Hub XMLRPC API service.
var HUB_XMLRPC_PORTS = []types.PortMap{
	NewPortMap(HubApiServiceName, "xmlrpc", 2830, 2830),
}

// PROXY_TCP_PORTS are the tcp ports required by the proxy.
var PROXY_TCP_PORTS = []types.PortMap{
	NewPortMap(ProxyTcpServiceName, "ssh", 8022, 22),
	NewPortMap(ProxyTcpServiceName, "publish", 4505, 4505),
	NewPortMap(ProxyTcpServiceName, "request", 4506, 4506),
}

// PROXY_PODMAN_PORTS are the http/s ports required by the proxy.
var PROXY_PODMAN_PORTS = []types.PortMap{
	NewPortMap(ProxyTcpServiceName, "https", 443, 443),
	NewPortMap(ProxyTcpServiceName, "http", 80, 80),
}

// GetProxyPorts returns all the proxy container ports.
func GetProxyPorts() []types.PortMap {
	ports := []types.PortMap{}
	ports = appendPorts(ports, false, PROXY_TCP_PORTS...)
	ports = appendPorts(ports, false, types.PortMap{
		Service:  ProxyUdpServiceName,
		Name:     "tftp",
		Exposed:  69,
		Port:     69,
		Protocol: "udp",
	})

	return ports
}
