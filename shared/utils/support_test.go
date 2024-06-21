// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"strings"
	"testing"
)

func TestGetSupportConfigPath(t *testing.T) {
	data := [][]string{
		{`/var/log/scc_uyuni-server.mgr.internal_240529_1124.txz`, `/var/log/scc_uyuni-server.mgr.internal_240529_1124.txz`},
		{`/var/log/scc_uyuni-server.mgr.internal.txz`, `/var/log/scc_uyuni-server.mgr.internal.txz`},
		{`/var/log/scc_uyuni-server_240529_1124.txz`, `/var/log/scc_uyuni-server_240529_1124.txz`},
		{`/var/log/scc_uyuni-server.txz`, `/var/log/scc_uyuni-server.txz`},
	}

	for i, testCase := range data {
		input := testCase[0]
		expected := testCase[1]

		actual := GetSupportConfigPath(input)

		if actual != expected {
			t.Errorf("Testcase %d: Expected %s got %s when GetSupportConfigPath %s", i, expected, actual, input)
		}
	}
}

func TestHostedContainers(t *testing.T) {
	data := `
    /etc/systemd/system/uyuni-server.service
	/etc/systemd/system/uyuni-server-attestation@.service
	`

	expected := []string{`uyuni-server`, `uyuni-server-attestation@`}

	actual := GetContainersFromSystemdFiles(data)

	if strings.Join(actual, " ") != strings.Join(expected, " ") {
		t.Errorf("Testcase: Expected %s got %s ", expected, actual)
	}
}
