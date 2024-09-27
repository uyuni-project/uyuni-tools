// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
)

func TestGetServerPorts(t *testing.T) {
	allPorts := len(WEB_PORTS) + len(PGSQL_PORTS) + len(SALT_PORTS) + len(COBBLER_PORTS) +
		len(TASKO_PORTS) + len(TOMCAT_PORTS) + len(SEARCH_PORTS) + len(TFTP_PORTS)

	ports := GetServerPorts(false)
	test_utils.AssertEquals(t, "Wrong number of ports without debug ones", allPorts-3, len(ports))

	ports = GetServerPorts(true)
	test_utils.AssertEquals(t, "Wrong number of ports with debug ones", allPorts, len(ports))
}
