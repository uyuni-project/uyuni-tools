// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import "github.com/uyuni-project/uyuni-tools/shared/types"

func NewPortMap(name string, exposed int, port int) types.PortMap {
	return types.PortMap{
		Name:    name,
		Exposed: exposed,
		Port:    port,
	}
}

// The port names should be less than 15 characters long and lowercased for traefik to eat them
var TCP_PORTS = []types.PortMap{
	NewPortMap("postgres", 5432, 5432),
	NewPortMap("salt-publish", 4505, 4505),
	NewPortMap("salt-request", 4506, 4506),
	NewPortMap("cobbler", 25151, 25151),
	NewPortMap("psql-mtrx", 9187, 9187),
	NewPortMap("tasko-jmx-mtrx", 5556, 5556),
	NewPortMap("tomcat-jmx-mtrx", 5557, 5557),
}

var DEBUG_PORTS = []types.PortMap{
	// We can't expose on port 8000 since traefik already uses it
	NewPortMap("tomcat-debug", 8003, 8003),
	NewPortMap("tasko-debug", 8001, 8001),
	NewPortMap("search-debug", 8002, 8002),
}

var UDP_PORTS = []types.PortMap{
	{
		Name:     "tftp",
		Exposed:  69,
		Port:     69,
		Protocol: "udp",
	},
}

var PROXY_TCP_PORTS = []types.PortMap{
	NewPortMap("ssh", 8022, 22),
	NewPortMap("https", 443, 443),
	NewPortMap("http", 80, 8080),
	NewPortMap("salt-publish", 4505, 4505),
	NewPortMap("salt-request", 4506, 4506),
}
