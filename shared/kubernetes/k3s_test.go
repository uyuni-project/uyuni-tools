// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"testing"
	"time"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Test that the generated endpoints are valid for traefik.
func TestGetTraefikEndpointName(t *testing.T) {
	ports := utils.GetServerPorts(true)
	ports = append(ports, utils.GetProxyPorts()...)

	for _, port := range ports {
		actual := GetTraefikEndpointName(port)
		// Traefik would fail if the name is longer than 15 characters
		if len(actual) > 15 {
			t.Errorf("Traefik endpoint name has more than 15 characters: %s", actual)
		}
	}
}

func TestWaitForTraefik(t *testing.T) {
	// Test that the time zone is properly handled
	installTime := time.Now().In(time.UTC).Add(time.Second * 42)
	newRunner = testutils.FakeRunnerGenerator(installTime.Format("2006-01-02T15:04:05Z"), nil)
	if err := waitForTraefik(); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}
