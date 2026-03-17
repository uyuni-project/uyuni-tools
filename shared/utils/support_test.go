// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"strings"
	"testing"
)

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
