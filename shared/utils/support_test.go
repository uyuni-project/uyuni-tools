// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestHostedContainers(t *testing.T) {
	data := `
    /etc/systemd/system/uyuni-server.service
	/etc/systemd/system/uyuni-server-attestation@.service
	`

	expected := []string{`uyuni-server`, `uyuni-server-attestation@`}

	actual := GetContainersFromSystemdFiles(data)

	testutils.AssertEquals(t, "Unexpected containers", expected, actual)
}
