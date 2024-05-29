// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
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
