// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package status

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestLogServerServicesState(t *testing.T) {
	cases := []struct {
		name     string
		enabled  []string
		expected string
	}{
		{name: "enabled", enabled: []string{podman.ServerService}, expected: "Server services are enabled"},
		{name: "disabled", enabled: []string{}, expected: "Server services are disabled"},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			driver := &testutils.FakeSystemdDriver{
				Installed: []string{podman.ServerService},
				Enabled:   testCase.enabled,
			}
			var output bytes.Buffer
			oldLogger := log.Logger
			log.Logger = zerolog.New(&output)
			t.Cleanup(func() {
				log.Logger = oldLogger
			})

			logServerServicesState(podman.NewSystemdWithDriver(driver))
			testutils.AssertStringContains(t, "service state missing: ", output.String(), testCase.expected)
		})
	}
}
