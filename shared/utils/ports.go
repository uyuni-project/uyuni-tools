// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import "github.com/uyuni-project/uyuni-tools/shared/types"

// NewPortMap is a constructor for PortMap type.
func NewPortMap(name string, exposed int, port int) types.PortMap {
	return types.PortMap{
		Name:    name,
		Exposed: exposed,
		Port:    port,
	}
}

// TCP_PORTS are the tcp ports required by the server
// The port names should be less than 15 characters long and lowercased for traefik to eat them.
var TCP_PORTS = []types.PortMap{
	NewPortMap("postgres", 5432, 5432),
	NewPortMap("salt-publish", 4505, 4505),
	NewPortMap("salt-request", 4506, 4506),
	NewPortMap("cobbler", 25151, 25151),
	NewPortMap("psql-mtrx", 9187, 9187),
	NewPortMap("tasko-jmx-mtrx", 5556, 5556),
	NewPortMap("tomcat-jmx-mtrx", 5557, 5557),
	NewPortMap("tasko-mtrx", 9800, 9800),
}

// TCP_PODMAN_PORTS are the tcp ports required by the server on podman.
var TCP_PODMAN_PORTS = []types.PortMap{
	// TODO: Replace Node exporter with cAdvisor
	NewPortMap("node-exporter", 9100, 9100),
}

// DEBUG_PORTS are the port used by dev for debugging applications.
var DEBUG_PORTS = []types.PortMap{
	// We can't expose on port 8000 since traefik already uses it
	NewPortMap("tomcat-debug", 8003, 8003),
	NewPortMap("tasko-debug", 8001, 8001),
	NewPortMap("search-debug", 8002, 8002),
}

// UDP_PORTS are the udp ports required by the server.
var UDP_PORTS = []types.PortMap{
	{
		Name:     "tftp",
		Exposed:  69,
		Port:     69,
		Protocol: "udp",
	},
}

// HUB_XMLRPC_PORTS are the tcp ports required by the Hub XMLRPC API service.
var HUB_XMLRPC_PORTS = []types.PortMap{
	NewPortMap("hub-xmlrpc", 2830, 2830),
}

// PROXY_TCP_PORTS are the tcp ports required by the proxy.
var PROXY_TCP_PORTS = []types.PortMap{
	NewPortMap("ssh", 8022, 22),
	NewPortMap("salt-publish", 4505, 4505),
	NewPortMap("salt-request", 4506, 4506),
}

// PROXY_PODMAN_PORTS are the http/s ports required by the proxy.
var PROXY_PODMAN_PORTS = []types.PortMap{
	NewPortMap("https", 443, 443),
	NewPortMap("http", 80, 80),
}
