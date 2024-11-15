// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestGetServiceImage(t *testing.T) {
	type dataType struct {
		catOut   string
		catErr   error
		expected string
	}
	data := []dataType{
		{"", errors.New("service not existing"), ""},
		{"content with no image defined", nil, ""},
		{`# /etc/systemd/system/uyuni-server-attestation@.service
[Unit]
Description=Uyuni server attestation container service
Wants=network.target
After=network-online.target
[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
[Install]
WantedBy=multi-user.target default.target

# /etc/systemd/system/uyuni-server-attestation@.service.d/generated.conf

[Service]
Environment=UYUNI_IMAGE=myregistry.org/silly/image:tag
`, nil, "myregistry.org/silly/image:tag"},
	}

	for _, testData := range data {
		runCmdOutput = func(_ zerolog.Level, _ string, _ ...string) ([]byte, error) {
			return []byte(testData.catOut), testData.catErr
		}

		testutils.AssertEquals(t, "Wrong image found", testData.expected, GetServiceImage("myservice"))
	}
}
