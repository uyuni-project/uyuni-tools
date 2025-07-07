// SPDX-FileCopyrightText: 2025 SUSE LLC
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
	// DBServiceName is the name of the server internal database service.
	DBServiceName = "db"
	// DBExporterServiceName is the name of the Prometheus database exporter service.
	DBExporterServiceName = "db"
	// TaskoServiceName is the name of the server taskomatic service.
	TaskoServiceName = "taskomatic"
	// TftpServiceName is the name of the server tftp service.
	TftpServiceName = "tftp"
	// TomcatServiceName is the name of the server tomcat service.
	TomcatServiceName = "tomcat"
	// SearchServiceName is the name of the server search service.
	SearchServiceName = "search"

	// HubAPIServiceName is the name of the server hub API service.
	HubAPIServiceName = "hub-api"

	// ProxyTCPServiceName is the name of the proxy TCP service.
	ProxyTCPServiceName = "uyuni-proxy-tcp"

	// ProxyUDPServiceName is the name of the proxy UDP service.
	ProxyUDPServiceName = "uyuni-proxy-udp"
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

// WebPorts is the list of ports for the server web service.
var WebPorts = []types.PortMap{
	NewPortMap(WebServiceName, "http", 80, 80),
}

// DBExporterPorts is the list of ports for the db exporter service.
var DBExporterPorts = []types.PortMap{
	NewPortMap(DBExporterServiceName, "exporter", 9187, 9187),
}

// ReportDBPorts is the list of ports for the server report db service.
var ReportDBPorts = []types.PortMap{
	NewPortMap(ReportdbServiceName, "pgsql", 5432, 5432),
}

// DBPorts is the list of ports for the server internal db service.
var DBPorts = []types.PortMap{
	NewPortMap(DBServiceName, "pgsql", 5432, 5432),
}

// SaltPorts is the list of ports for the server salt service.
var SaltPorts = []types.PortMap{
	NewPortMap(SaltServiceName, "publish", 4505, 4505),
	NewPortMap(SaltServiceName, "request", 4506, 4506),
}

// CobblerPorts is the list of ports for the server cobbler service.
var CobblerPorts = []types.PortMap{
	NewPortMap(CobblerServiceName, "cobbler", 25151, 25151),
}

// TaskoPorts is the list of ports for the server taskomatic service.
var TaskoPorts = []types.PortMap{
	NewPortMap(TaskoServiceName, "jmx", 5556, 5556),
	NewPortMap(TaskoServiceName, "mtrx", 9800, 9800),
	NewPortMap(TaskoServiceName, "debug", 8001, 8001),
}

// TomcatPorts is the list of ports for the server tomcat service.
var TomcatPorts = []types.PortMap{
	NewPortMap(TomcatServiceName, "jmx", 5557, 5557),
	NewPortMap(TomcatServiceName, "debug", 8003, 8003),
}

// SearchPorts is the list of ports for the server search service.
var SearchPorts = []types.PortMap{
	NewPortMap(SearchServiceName, "debug", 8002, 8002),
}

// TftpPorts is the list of ports for the server tftp service.
var TftpPorts = []types.PortMap{
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
	ports = appendPorts(ports, debug, WebPorts...)
	ports = appendPorts(ports, debug, SaltPorts...)
	ports = appendPorts(ports, debug, CobblerPorts...)
	ports = appendPorts(ports, debug, TaskoPorts...)
	ports = appendPorts(ports, debug, TomcatPorts...)
	ports = appendPorts(ports, debug, SearchPorts...)
	ports = appendPorts(ports, debug, TftpPorts...)
	ports = appendPorts(ports, debug, DBExporterPorts...)

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
